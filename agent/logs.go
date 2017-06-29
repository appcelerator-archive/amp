package core

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/appcelerator/amp/api/rpc/logs"
	"github.com/appcelerator/amp/pkg/nats-streaming"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/gogo/protobuf/proto"
)

const (
	stdWriterPrefixLen = 8
	stdWriterFdIndex   = 0
	stdWriterSizeIndex = 4

	startingBufLen = 32*1024 + stdWriterPrefixLen + 1
)

// verify all containers to open logs stream if not already done
func (a *Agent) updateLogsStream() {
	for ID, data := range a.containers {
		if data.logsStream == nil || data.logsReadError {
			lastTimeID := a.getLastTimeID(ID)
			if lastTimeID == "" {
				log.Infof("open logs stream from the beginning container %s\n", data.name)
			} else {
				log.Infof("open logs stream from time_id=%s on container %s\n", lastTimeID, data.name)
			}
			stream, err := a.openLogsStream(ID, lastTimeID)
			if err != nil {
				log.Errorf("Error opening logs stream on container: %s\n", data.name)
			} else {
				data.logsStream = stream

				// Inspect the container to check if it's a TTY
				c, err := a.dock.GetClient().ContainerInspect(context.Background(), ID)
				if err != nil {
					log.Errorf("Error inspecting container for TTY: %s\n", data.name)
					continue
				}

				// Read logs in the background
				go func(ID string, data *ContainerData, tty bool) {
					// Pick the adequate log reader function
					logReader := a.readLogs
					if tty {
						logReader = a.readLogsTTY
					}

					// Read logs
					log.Infof("start reading log stream of container: %s\n", data.name)
					if err := logReader(ID, data); err != nil {
						log.Errorf("Error reading log stream of container %s: %v", data.name, err)
					}
					log.Infof("stop reading log stream of container: %s\n", data.name)

					// Close log stream
					if err := data.logsStream.Close(); err != nil {
						log.Errorf("Error closing log stream of container %s: %v", data.name, err)
					}
				}(ID, data, c.Config.Tty)
			}
		}
	}
}

// open a logs container stream
func (a *Agent) openLogsStream(ID string, lastTimeID string) (io.ReadCloser, error) {
	containerLogsOptions := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Timestamps: true,
	}
	if lastTimeID != "" {
		containerLogsOptions.Since = lastTimeID
	}
	return a.dock.GetClient().ContainerLogs(context.Background(), ID, containerLogsOptions)
}

// get last timestamp if exist
func (a *Agent) getLastTimeID(ID string) string {
	data, err := ioutil.ReadFile(path.Join(containersDataDir, ID))
	if err != nil {
		return ""
	}
	return string(data)
}

func (a *Agent) buildLogEntry(ID string, data *ContainerData, line string, timeID int64) *logs.LogEntry {
	// Log entry is formatted with 1/ a 30 characters timestamp, 2/ the message content, for instance:
	// 2017-06-14T20:07:57.425972267Z building a seeds list for cluster etcd
	date := line[:30]
	msg := strings.TrimSuffix(line[31:], "\n")

	// Convert date to timestamp
	timestamp, err := time.Parse("2006-01-02T15:04:05.999999999Z", date)
	if err != nil {
		timestamp = time.Now()
	}

	// Build log entry
	return &logs.LogEntry{
		Timestamp:          timestamp.Format(time.RFC3339Nano),
		ContainerId:        ID,
		ContainerName:      data.name,
		ContainerShortName: data.shortName,
		ContainerState:     data.state,
		ServiceName:        data.serviceName,
		ServiceId:          data.serviceID,
		TaskId:             data.taskID,
		TaskSlot:           int32(data.taskSlot),
		StackName:          data.stackName,
		NodeId:             data.nodeID,
		TimeId:             fmt.Sprintf("%016X", timeID),
		Labels:             data.labels,
		Msg:                msg,
	}
}

// readLogsTTY read logs from a TTY container.
func (a *Agent) readLogsTTY(ID string, data *ContainerData) error {
	var (
		previous, now int64
		br            = bufio.NewReader(data.logsStream)
	)

	for {
		line, err := br.ReadString('\n')
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		// Compute TimeId
		now = time.Now().UnixNano()
		if now <= previous {
			now = previous + 1
		}
		previous = now

		logEntry := a.buildLogEntry(ID, data, line, now)
		a.addLogEntry(logEntry, data, logEntry.Timestamp)
	}
}

// readLogs is a modified version of docker stdcopy.StdCopy.
func (a *Agent) readLogs(ID string, data *ContainerData) error {
	var (
		buf           = make([]byte, startingBufLen)
		bufLen        = len(buf)
		src           = data.logsStream
		nr            int
		frameSize     int
		previous, now int64
	)

	for {
		// Make sure we have at least a full header
		for nr < stdWriterPrefixLen {
			nr2, err := src.Read(buf[nr:])
			nr += nr2
			if err == io.EOF {
				if nr < stdWriterPrefixLen {
					return nil
				}
				break
			}
			if err != nil {
				return err
			}
		}

		stream := stdcopy.StdType(buf[stdWriterFdIndex])
		// Check the first byte to know where to write
		switch stream {
		case stdcopy.Stdin:
		case stdcopy.Stdout:
		case stdcopy.Stderr:
		case stdcopy.Systemerr:
		default:
			return fmt.Errorf("Unrecognized input header: %d", buf[stdWriterFdIndex])
		}

		// Retrieve the size of the frame
		frameSize = int(binary.BigEndian.Uint32(buf[stdWriterSizeIndex : stdWriterSizeIndex+4]))

		// Check if the buffer is big enough to read the frame.
		// Extend it if necessary.
		if frameSize+stdWriterPrefixLen > bufLen {
			buf = append(buf, make([]byte, frameSize+stdWriterPrefixLen-bufLen+1)...)
			bufLen = len(buf)
		}

		// While the amount of bytes read is less than the size of the frame + header, we keep reading
		for nr < frameSize+stdWriterPrefixLen {
			nr2, err := src.Read(buf[nr:])
			nr += nr2
			if err == io.EOF {
				if nr < frameSize+stdWriterPrefixLen {
					return nil
				}
				break
			}
			if err != nil {
				return err
			}
		}

		// we might have an error from the source mixed up in our multiplexed
		// stream. if we do, return it.
		if stream == stdcopy.Systemerr {
			return fmt.Errorf("error from daemon in stream: %s", string(buf[stdWriterPrefixLen:frameSize+stdWriterPrefixLen]))
		}

		// Compute TimeId
		now = time.Now().UnixNano()
		if now <= previous {
			now = previous + 1
		}
		previous = now

		// line contains a full log entry
		line := string(buf[stdWriterPrefixLen : frameSize+stdWriterPrefixLen])

		logEntry := a.buildLogEntry(ID, data, line, now)
		a.addLogEntry(logEntry, data, logEntry.Timestamp)

		// Move the rest of the buffer to the beginning
		copy(buf, buf[frameSize+stdWriterPrefixLen:])

		// Move the index
		nr -= frameSize + stdWriterPrefixLen
	}
}

func (a *Agent) addLogEntry(entry *logs.LogEntry, data *ContainerData, date string) {
	if conf.logsBufferPeriod == 0 || conf.logsBufferSize == 0 {
		a.logsBuffer.Entries[0] = entry
		a.sendLogsBuffer()
		a.periodicDataSave(data, date)
		return
	}
	a.logsBufferMutex.Lock()
	defer a.logsBufferMutex.Unlock()
	if a.logsBuffer == nil {
		a.logsBuffer.Entries = make([]*logs.LogEntry, conf.logsBufferSize)
	}
	a.logsBuffer.Entries = append(a.logsBuffer.Entries, entry)
	if len(a.logsBuffer.Entries) >= conf.logsBufferSize {
		a.sendLogsBuffer()
		a.logsBuffer.Entries = nil
		a.periodicDataSave(data, date)
	}
}

func (a *Agent) sendLogsBuffer() {
	encoded, err := proto.Marshal(a.logsBuffer)
	if err != nil {
		log.Errorf("error marshalling log entries: %v\n", err)
		return
	}
	_, err = a.natsStreaming.GetClient().PublishAsync(ns.LogsSubject, encoded, nil)
	if err != nil {
		log.Errorf("error sending log entry: %v\n", err)
		return
	}
	a.nbLogs += len(a.logsBuffer.Entries)
}

func (a *Agent) periodicDataSave(data *ContainerData, date string) {
	now := time.Now()
	if now.Sub(data.lastDateSaveTime).Seconds() >= float64(a.logsSavedDatePeriod) {
		err := ioutil.WriteFile(path.Join(containersDataDir, data.ID), []byte(date), 0666)
		if err != nil {
			log.Errorf("error writing to container data directory: ", err)
		}
		data.lastDateSaveTime = now
	}
}

// close all logs stream
func (a *Agent) closeLogsStreams() {
	for _, data := range a.containers {
		if data.logsStream != nil {
			err := data.logsStream.Close()
			if err != nil {
				log.Errorf("Error closing a log stream: ", err)
			}
		}
	}
}

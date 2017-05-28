
// copy amp stats proto file
export class StatsRequest {
  stats_cpu: boolean
  stats_mem: boolean
  stats_io: boolean
  stats_net: boolean
  group: string
  filter_container_id: string
  filter_container_name: string
  filter_container_short_name: string
  filter_container_state: string
  filter_service_name: string
  filter_service_id: string
  filter_task_id: string
  filter_stack_name: string
  filter_node_id: string
  period: string
  time_group: string
  avg: boolean
  //
  format: boolean

  constructor() {
  }

}

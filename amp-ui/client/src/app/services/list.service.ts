
export class ListService {
  allData : any[] = []
  filteredData : any[] = []
  filteredDataSize = 0
  filterValue = ""
  filterFunction : (any, string) => boolean
  pageIndex = 0
  pageSize = 0

  constructor() {}

  setFilterFunction(filter : (any, string) => boolean) {
    this.filterFunction = filter
  }
  setData(data : any[]) {
    this.allData = data
    this.filter(this.filterValue)
  }

  getData() {
    return this.filteredData
  }

  filter(value : string) {
    this.filterValue = value
    if (value=='') {
      this.filteredData = this.page(this.allData.slice());
    }
    this.filteredData = this.page(
        this.allData.filter(
          (item) => {
            return this.filterFunction(item, value)
          }
        )
      );
  }

  order(field : string, asc : number) {
    this.filteredData.sort(
      (a, b) => {
        return a[field].localeCompare(b[field]) * asc;
      }
    );
  }

  orderNum(field : string, asc : number) {
    this.filteredData.sort(
      (a, b) => {
        return (a[field]-b[field]) * asc;
      }
    );
  }

  page(list : any[]) : any[] {
    this.filteredDataSize = list.length
    if (this.pageSize == 0) {
      return list
    }
    if (this.pageIndex * this.pageSize >= list.length) {
      return []
    }
    if (this.pageIndex * this.pageSize + this.pageSize >= list.length) {
      return list.slice(this.pageIndex * this.pageSize)
    }
    return list.slice(this.pageIndex * this.pageSize, this.pageIndex * this.pageSize + this.pageSize)
  }

  setPageSize( size : number) {
    this.pageSize = size
    this.filter(this.filterValue)
  }

  setPageIndex(relative : number, absolute : number) {
    //console.log("r="+relative+", abs="+absolute+", index="+this.pageIndex)
    if (relative == 0) {
        this.pageIndex = absolute - 1
    }  else if (relative == 1) {
        this.pageIndex++
    } else if (relative == -1) {
        this.pageIndex--
    }
    //console.log("dataLength="+this.filteredDataSize+", pageSize="+this.pageSize)
    let index=Math.floor(this.filteredDataSize / this.pageSize)
    let modulus=this.filteredDataSize % this.pageSize
    //console.log("newindexPage="+this.pageIndex+", index="+index+", modulo="+modulus)
    if (this.pageIndex > index) {
      this.pageIndex = index
    } else if (this.pageIndex == index && modulus == 0) {
      this.pageIndex = index - 1
    }
    if (this.pageIndex<0) {
      this.pageIndex = 0
    }
    //console.log("index="+this.pageIndex)
    this.filter(this.filterValue)
  }

}

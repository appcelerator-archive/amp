
export class ListService {
  allData : any[] = []
  filteredData : any[] = []
  filterValue = ""
  filterFunction : (any, string) => boolean

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
      this.filteredData = this.allData.slice();
    }
    this.filteredData = this.allData.filter(
      (item) => {
        return this.filterFunction(item, value)
      }
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

}

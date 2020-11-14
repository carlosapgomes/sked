import React, { Component } from "react";
import dayjs from "dayjs";

export default class ScheduleList extends Component {
  constructor(props) {
    super(props);

    this.state = {
      currentMonth: "",
      days: undefined,
    };
  }
  componentDidMount() {
    let m = dayjs().format("YYYY-MM");
    this.setState({
      currentMonth: m,
    });
    const nDays = new Date(m.substr(0, 4), m.substr(5, 2), 0).getDate();
    let days = [];
    for (let i = 1; i <= nDays; i++) {
      days.push(<li key={i.toString()}>{i}</li>);
    }
    this.setState({
      days: [...days],
    });
  }
  setCurrentMonth(m) {
    this.setState({
      currentMonth: m,
    });
    const nDays = new Date(m.substr(0, 4), m.substr(5, 2), 0).getDate();
    let days = [];
    for (let i = 1; i <= nDays; i++) {
      days.push(<li key={i.toString()}>{i}</li>);
    }
    this.setState({
      days: [...days],
    });
  }
  render() {
    return (
      <div>
        <input
          type="month"
          value={this.state.currentMonth}
          name="month"
          id="month"
          onChange={(e) => {
            this.setCurrentMonth(e.target.value);
          }}
        />
        <div>
          <ul>{this.state.days}</ul>
        </div>
      </div>
    );
  }
}

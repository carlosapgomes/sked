import React, { Component } from "react";
import dayjs from "dayjs";
import cl from "./ScheduleList.css";
export default class ScheduleList extends Component {
  constructor(props) {
    super(props);

    this.state = {
      currentMonth: "",
      days: undefined,
      appointments: undefined,
      surgeries: undefined,
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
  getAllAppointmentsInAMonth(m, y) {
    let ajax = new XMLHttpRequest();
    let url = "https://dev.local/api/appointments?month=" + m + "&year=" + y;
    ajax.open("GET", url, true);
    ajax.withCredentials = true;
    ajax.setRequestHeader("Content-type", "application/json");
    ajax.send();
    ajax.onreadystatechange = () => {
      if (ajax.readyState === 4 && ajax.status === 200) {
        if (!ajax.responseText) {
          window.alert("Could not find any patient");
        } else {
          let data = JSON.parse(ajax.responseText);
          if (data) {
            console.log(data);
            this.setState({
              appointments: [...data],
            });
            this.updateSchedulesList(m.y);
          }
        }
      }
      if (ajax.readyState === 4 && ajax.status !== 200) {
        window.alert("Could not complete the operation");
        console.log(ajax.responseText);
      }
    };
  }
  getAllSurgeriesInAMonth(m, y) {
    let ajax = new XMLHttpRequest();
    let url = "https://dev.local/api/surgeries?month=" + m + "&year=" + y;
    ajax.open("GET", url, true);
    ajax.withCredentials = true;
    ajax.setRequestHeader("Content-type", "application/json");
    ajax.send();
    ajax.onreadystatechange = () => {
      if (ajax.readyState === 4 && ajax.status === 200) {
        if (!ajax.responseText) {
          window.alert("Could not find any patient");
        } else {
          let data = JSON.parse(ajax.responseText);
          if (data) {
            console.log(data);
            this.setState({
              surgeries: [...data],
            });
            this.updateSchedulesList(m, y);
          }
        }
      }
      if (ajax.readyState === 4 && ajax.status !== 200) {
        window.alert("Could not complete the operation");
        console.log(ajax.responseText);
      }
    };
  }
  updateSchedulesList(m, y) {
    if (
      typeof this.state.appointments === "undefined" ||
      typeof this.state.surgeries === "undefined"
    ) {
      return;
    }
    const nDays = new Date(y, m, 0).getDate();
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
      appointments: undefined,
      surgeries: undefined,
    });
    let year = m.substr(0, 4);
    let month = m.substr(5, 2);
    this.getAllSurgeriesInAMonth(month, year);
    this.getAllAppointmentsInAMonth(month, year);
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
          <ul className={cl.ul}>{this.state.days}</ul>
        </div>
      </div>
    );
  }
}

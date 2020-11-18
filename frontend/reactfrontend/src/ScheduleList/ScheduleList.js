import React, { Component } from "react";
import dayjs from "dayjs";
import cl from "./ScheduleList.css";
export default class ScheduleList extends Component {
  constructor(props) {
    super(props);

    this.state = {
      currentMonth: "",
      currentYear: "",
      days: undefined,
      appointments: undefined,
      surgeries: undefined,
    };
  }
  componentDidMount() {
    const m = String(dayjs().month());
    const y = String(dayjs().year());
    this.setState({
      currentMonth: m,
      currentYear: y,
    });
    //const nDays = new Date(m.substr(0, 4), m.substr(5, 2), 0).getDate();
    //let days = [];
    //for (let i = 1; i <= nDays; i++) {
    //days.push(<li key={i.toString()}>{i}</li>);
    //}
    //this.setState({
    //days: [...days],
    //});
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
    //const nDays = new Date(y, m, 0).getDate();
    let daysOfMonth = [];
    //for (let i = 1; i <= nDays; i++) {
    //daysOfMonth.push(
    //<li key={i.toString()}>
    //{this.state.appointments.map((e) => {
    //if (dayjs(e.date_time).date() === i) {
    //return (
    //<div>
    //{e.doctor_name}: {e.patient_name}
    //</div>
    //);
    //}
    //})}
    //</li>
    //);
    //}
    this.setState({
      days: [...daysOfMonth],
    });
  }
  setCurrentMonth(m) {
    this.setState({
      currentMonth: m,
      appointments: undefined,
      surgeries: undefined,
    });
    console.log(this.state.currentMonth);
    console.log(this.state.currentYear);

    //let year = m.substr(0, 4);
    //let month = m.substr(5, 2);
    //this.getAllSurgeriesInAMonth(month, year);
    //this.getAllAppointmentsInAMonth(month, year);
  }
  setCurrentYear(y) {
    this.setState({
      currentYear: y,
      appointments: undefined,
      surgeries: undefined,
    });
    console.log(this.state.currentMonth);
    console.log(this.state.currentYear);
  }
  render() {
    return (
      <div>
        <span>
          <label for="month">Month: </label>
          <select
            id="month"
            name="month"
            value={this.state.currentMonth}
            onChange={(e) => {
              this.setCurrentMonth(e.target.value);
            }}
          >
            <option value="01">01</option>
            <option value="02">02</option>
            <option value="03">03</option>
            <option value="04">04</option>
            <option value="05">05</option>
            <option value="06">06</option>
            <option value="07">07</option>
            <option value="08">08</option>
            <option value="09">09</option>
            <option value="10">10</option>
            <option value="11">11</option>
            <option value="12">12</option>
          </select>
        </span>
        {"    "}
        {"    "}
        <span>
          <label for="year">Year: </label>
          <select
            id="year"
            name="year"
            value={this.state.currentYear}
            onChange={(e) => {
              this.setCurrentYear(e.target.value);
            }}
          >
            <option value="2020">2020</option>
            <option value="2021">2021</option>
          </select>
        </span>
        <div>
          <ul className={cl.ul}>{this.state.days}</ul>
        </div>
      </div>
    );
  }
}

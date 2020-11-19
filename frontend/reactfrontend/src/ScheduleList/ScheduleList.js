import React, { Component } from "react";
import dayjs from "dayjs";
import cl from "./ScheduleList.css";
export default class ScheduleList extends Component {
  constructor(props) {
    super(props);

    this.state = {
      currentMonth: "",
      currentYear: "",
      appointments: [],
      appSchedules: [],
      appSelected: false,
      surgeries: [],
      surgSchedules: [],
    };
  }
  componentDidMount() {
    const m = dayjs().month();
    const y = dayjs().year();
    this.setState({
      currentMonth: String(m),
      currentYear: String(y),
      appSelected: true,
    });
    this.updateAppointmtsAndSurgsData(m, y);
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
            this.setState({
              appointments: [...data],
            });
          }
          this.updateAppSchedulesList(m, y, data);
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
            this.setState({
              surgeries: [...data],
            });
          }
          this.updateSurgSchedulesList(m, y, data);
        }
      }
      if (ajax.readyState === 4 && ajax.status !== 200) {
        window.alert("Could not complete the operation");
        console.log(ajax.responseText);
      }
    };
  }
  updateAppSchedulesList(m, y, data) {
    const nDays = new Date(y, m, 0).getDate();
    let daysOfMonth = [];
    for (let i = 1; i <= nDays; i++) {
      daysOfMonth.push(
        <li key={i.toString()}>
          <div>
            {
              <span>
                {i}/{m}/{y}:{" "}
              </span>
            }
            <div>
              Appointments:
              {!data
                ? null
                : data.map((e) => {
                    let d = dayjs(e.dateTime).date();
                    if (d === i) {
                      return (
                        <div key={e.id}>
                          {e.doctorName} : {e.patientName}
                        </div>
                      );
                    } else {
                      return <div key={e.id}>{"  "}</div>;
                    }
                  })}
            </div>
          </div>
        </li>
      );
    }
    this.setState({
      appSchedules: [...daysOfMonth],
    });
  }
  updateSurgSchedulesList(m, y, data) {
    const nDays = new Date(y, m, 0).getDate();
    let daysOfMonth = [];
    for (let i = 1; i <= nDays; i++) {
      daysOfMonth.push(
        <li key={i.toString()}>
          <div>
            {
              <span>
                {i}/{m}/{y}:{" "}
              </span>
            }
            <div>
              Surgeries:
              {!data
                ? null
                : data.map((e) => {
                    let d = dayjs(e.dateTime).date();
                    if (d === i) {
                      return (
                        <div key={e.id}>
                          {e.doctorName} : {e.patientName}
                        </div>
                      );
                    } else {
                      return <div key={e.id}>{"  "}</div>;
                    }
                  })}
            </div>
          </div>
        </li>
      );
    }
    this.setState({
      surgSchedules: [...daysOfMonth],
    });
  }
  setCurrentMonth(m) {
    this.setState({
      currentMonth: m,
      appointments: [],
      surgeries: [],
    });
    this.updateAppointmtsAndSurgsData(
      Number(m),
      Number(this.state.currentYear)
    );
  }
  updateAppointmtsAndSurgsData(m, y) {
    this.getAllSurgeriesInAMonth(m, y);
    this.getAllAppointmentsInAMonth(m, y);
  }
  setCurrentYear(y) {
    this.setState({
      currentYear: y,
      appointments: [],
      surgeries: [],
    });
    this.updateAppointmtsAndSurgsData(
      Number(this.state.currentMonth),
      Number(y)
    );
  }
  radioChanged(e) {
    if (e.target.value === "appointments") {
      this.setState({
        appSelected: true,
      });
    } else {
      this.setState({
        appSelected: false,
      });
    }
  }
  render() {
    return (
      <div>
        <div>
          <input
            type="radio"
            value="appointments"
            name="schedule"
            id="appointments"
            checked={this.state.appSelected}
            onChange={(e) => {
              this.radioChanged(e);
            }}
          />
          <label htmlFor="appointments">Appointments</label>
          <input
            type="radio"
            value="surgeries"
            name="schedule"
            id="surgeries"
            checked={!this.state.appSelected}
            onChange={(e) => {
              this.radioChanged(e);
            }}
          />
          <label htmlFor="surgeries">Surgeries</label>
        </div>
        <span>
          <label htmlFor="month">Month: </label>
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
          <label htmlFor="year">Year: </label>
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
          <ul className={cl.ul}>
            {this.state.appSelected
              ? this.state.appSchedules
              : this.state.surgSchedules}
          </ul>
        </div>
      </div>
    );
  }
}

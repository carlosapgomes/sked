import React, { Component } from "react";
import dayjs from "dayjs";
import weekday from "dayjs/plugin/weekday";
import "./ScheduleList.css";
dayjs.extend(weekday);
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
    let month = Number(m) + 1;
    let ajax = new XMLHttpRequest();
    let url =
      "https://dev.local/api/appointments?month=" +
      String(month) +
      "&year=" +
      y;
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
    let month = Number(m) + 1;
    let ajax = new XMLHttpRequest();
    let url =
      "https://dev.local/api/surgeries?month=" + String(month) + "&year=" + y;
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
  clickedOnDay(d) {
    console.log(d);
  }
  clickedOnAppt(a) {
    console.log(a);
  }
  clickedOnSurg(s) {
    console.log(s);
  }
  updateAppSchedulesList(m, y, data) {
    // next month day zero corresponds to last day of the current month
    const nDays = new Date(y, Number(m) + 1, 0).getDate();
    let daysOfMonth = [];
    for (let i = 0; i < nDays; i++) {
      let nOfSchedules = 0;
      let schedules = [];
      let thisDay = dayjs(new Date(y, m, i + 1));
      if (data) {
        data.forEach((e) => {
          let day = dayjs(e.dateTime).date();
          if (day === i + 1) {
            nOfSchedules++;
            schedules.push(e);
          }
        });
      }
      daysOfMonth.push({
        dateStr: thisDay.format("DD/MM/YYYY"),
        nOfSchedules: nOfSchedules,
        schedules: [...schedules],
        weekday: thisDay.format("dd"),
      });
    }
    let apptsSchedules = [];
    for (let i = 0; i < nDays; i++) {
      apptsSchedules.push(
        <li key={(i + 1).toString()}>
          <details>
            {
              <summary>
                {daysOfMonth[i].dateStr}: {daysOfMonth[i].weekday} - #&nbsp;
                {daysOfMonth[i].nOfSchedules}
              </summary>
            }
            <div>
              <span
                className="AddSchedule"
                role="img"
                aria-label="Add appointment"
                data-day={i + 1}
                onClick={(e) => {
                  this.clickedOnDay(e.target.dataset.day);
                }}
              >
                &#10133;
              </span>
              {!daysOfMonth[i].schedules
                ? null
                : daysOfMonth[i].schedules.map((e) => {
                    console.log(e);
                    return (
                      <div
                        key={e.id}
                        data-id={e.id}
                        onClick={(e) => {
                          this.clickedOnAppt(e.target.dataset.id);
                        }}
                      >
                        {e.doctorName} : {e.patientName}
                      </div>
                    );
                  })}
            </div>
          </details>
        </li>
      );
    }
    this.setState({
      appSchedules: [...apptsSchedules],
    });
  }
  updateSurgSchedulesList(m, y, data) {
    const nDays = new Date(y, Number(m) + 1, 0).getDate();
    let daysOfMonth = [];
    for (let i = 0; i < nDays; i++) {
      let nOfSchedules = 0;
      let schedules = [];
      let thisDay = dayjs(new Date(y, m, i + 1));
      if (data) {
        data.forEach((e) => {
          let day = dayjs(e.dateTime).date();
          if (day === i + 1) {
            nOfSchedules++;
            schedules.push(e);
          }
        });
      }
      daysOfMonth.push({
        dateStr: thisDay.format("DD/MM/YYYY"),
        nOfSchedules: nOfSchedules,
        schedules: [...schedules],
        weekday: thisDay.format("dd"),
      });
    }
    let surgsSchedules = [];
    for (let i = 0; i < nDays; i++) {
      surgsSchedules.push(
        <li key={(i + 1).toString()}>
          <details>
            {
              <summary>
                {daysOfMonth[i].dateStr}: {daysOfMonth[i].weekday} - #&nbsp;
                {daysOfMonth[i].nOfSchedules}
              </summary>
            }
            <div>
              <span
                className="AddSchedule"
                role="img"
                aria-label="Add appointment"
                data-day={i}
                onClick={(e) => {
                  this.clickedOnDay(e.target.dataset.day);
                }}
              >
                &#10133;
              </span>
              {!daysOfMonth[i].schedules
                ? null
                : daysOfMonth[i].schedules.map((e) => {
                    return (
                      <div
                        key={e.id}
                        data-id={e.id}
                        onClick={(e) => {
                          this.clickedOnAppt(e.target.dataset.id);
                        }}
                      >
                        {e.doctorName} : {e.patientName}
                        <br />
                        {e.proposedSurgery}
                      </div>
                    );
                  })}
            </div>
          </details>
        </li>
      );
    }
    this.setState({
      surgSchedules: [...surgsSchedules],
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
        <h1>Schedules</h1>
        <section>
          <form>
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
            <br />
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
            <hr />
            <p>Choose Month/Year:</p>
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
                <option value="0">01</option>
                <option value="01">02</option>
                <option value="02">03</option>
                <option value="03">04</option>
                <option value="04">05</option>
                <option value="05">06</option>
                <option value="06">07</option>
                <option value="07">08</option>
                <option value="08">09</option>
                <option value="09">10</option>
                <option value="10">11</option>
                <option value="11">12</option>
              </select>
              {"    "}
              {"    "}
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
            <hr />
            <p>Days:</p>
            <div className="DaysList">
              <ul>
                {this.state.appSelected
                  ? this.state.appSchedules
                  : this.state.surgSchedules}
              </ul>
            </div>
          </form>
        </section>
      </div>
    );
  }
}

import React, { Component } from "react";
import dayjs from "dayjs";
import weekday from "dayjs/plugin/weekday";
import "./ScheduleList.css";
import { withTranslation } from "react-i18next";
import calendarize from "calendarize";
import Calendar from "../Calendar/Calendar";
dayjs.extend(weekday);

class ScheduleList extends Component {
  constructor(props) {
    super(props);

    this.state = {
      lang: "",
      currentMonth: "",
      currentYear: "",
      appointments: [],
      parsedAppointments: [],
      appSelected: false,
      surgeries: [],
      parsedSurgeries: [],
      calendarView: [],
      schedsInSelectedDay: [],
      selectedDay: "",
    };
    this.myRef = React.createRef();
  }
  componentDidMount() {
    const lang = this.props.i18n.language.toString().toLowerCase();
    (async () => {
      if (lang === "pt-br") {
        await import("dayjs/locale/pt-br.js");
      }
      dayjs.locale(lang);
      this.setState({
        lang: lang,
      });
    })();
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
    let url = "/api/appointments?month=" + String(month) + "&year=" + y;
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
    let url = "/api/surgeries?month=" + String(month) + "&year=" + y;
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
  selectCalendarDay(d) {
    this.setState({
      selectedDay: d,
    });
    let schedsInSelectedDay = [];
    if (this.state.appSelected) {
      // show appointments for day d
      this.state.parsedAppointments.forEach((e) => {
        if (e.day === d) {
          e.schedules.forEach((s) => {
            schedsInSelectedDay.push({
              id: s.id,
              doctorName: s.doctorName,
              patientName: s.patientName,
            });
          });
        }
      });
    } else {
      // show surgeries for day d
      this.state.parsedSurgeries.forEach((e) => {
        if (e.day === d) {
          e.schedules.forEach((s) => {
            schedsInSelectedDay.push({
              id: s.id,
              doctorName: s.doctorName,
              patientName: s.patientName,
            });
          });
        }
      });
    }
    this.setState({
      schedsInSelectedDay: [...schedsInSelectedDay],
    });
    if (window.screen.width < 768) {
      this.myRef.current.scrollIntoView({ behavior: "smooth" });
    }
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
        day: thisDay.format("D"),
      });
    }
    this.setState({
      parsedAppointments: [...daysOfMonth],
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
        day: thisDay.format("D"),
      });
    }
    this.setState({
      parsedSurgeries: [...daysOfMonth],
    });
  }
  setCurrentMonth(m) {
    this.setState({
      currentMonth: m,
      appointments: [],
      parsedAppointments: [],
      surgeries: [],
      parsedSurgeries: [],
    });
    this.updateAppointmtsAndSurgsData(
      Number(m),
      Number(this.state.currentYear)
    );
  }
  updateAppointmtsAndSurgsData(m, y) {
    this.getAllSurgeriesInAMonth(m, y);
    this.getAllAppointmentsInAMonth(m, y);
    const date = new Date(Number(y), Number(m), 1);
    let view = calendarize(date);
    this.setState({ calendarView: [...view] });
  }
  setCurrentYear(y) {
    this.setState({
      currentYear: y,
      appointments: [],
      parsedAppointments: [],
      surgeries: [],
      parsedSurgeries: [],
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
    const { t } = this.props;
    return (
      <div className="SchedulesList">
        <section style={{ display: "flex", flexDirection: "column" }}>
          <form>
            <div style={{ display: "flex", justifyContent: "space-evenly" }}>
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
              <label htmlFor="appointments">{t("Appointments")}</label>
              &nbsp;&nbsp;
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
              <label htmlFor="surgeries">{t("Surgeries")}</label>
            </div>{" "}
            <div style={{ display: "flex", justifyContent: "space-evenly" }}>
              <label htmlFor="month">{t("Month")}: </label>
              &nbsp;
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
              &nbsp;&nbsp;
              <label htmlFor="year">{t("Year")}: </label>
              &nbsp;
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
            </div>
            <Calendar
              view={this.state.calendarView}
              selectCalendarDay={(d) => this.selectCalendarDay(d)}
            ></Calendar>
          </form>
        </section>
        <section className="DayListSection">
          <p>
            <b ref={this.myRef}>
              {t("Schedule")}: {t("Day")} - {this.state.selectedDay}
            </b>
          </p>
          <div>
            <ul>
              {this.state.schedsInSelectedDay.map((e) => {
                return (
                  <li key={e.id}>
                    {e.doctorName}: {e.patientName}
                  </li>
                );
              })}
            </ul>
          </div>
        </section>
      </div>
    );
  }
}

export default withTranslation()(ScheduleList);

import React, { Component } from "react";
import { BrowserRouter as Router, Switch, Route, Link } from "react-router-dom";
import "./App.css";
import Auth from "./Auth/Auth";
import Appointments from "./Appointments/Appointments";
import Surgeries from "./Surgeries/Surgeries";
import Patients from "./Patients/Patients";
import Users from "./Users/Users";
import ScheduleList from "./ScheduleList/ScheduleList";
class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      currentUser: null,
      loggedIn: false,
      currentPatient: null,
      currentDoctor: null,
      doctors: [],
    };
  }
  componentDidMount() {
    // try to recover user data in case of a page reload
    if (window.localStorage.getItem("email") !== "") {
      let ajax = new XMLHttpRequest();
      let url =
        "https://dev.local/api/users?email=" +
        window.localStorage.getItem("email");
      ajax.open("GET", url, true);
      ajax.withCredentials = true;
      ajax.setRequestHeader("Content-type", "application/json");
      ajax.send();
      ajax.onreadystatechange = () => {
        if (ajax.readyState === 4 && ajax.status === 200) {
          if (!ajax.responseText) {
            console.log("Received an empty user");
          }
          // update state & localStorage
          let data = JSON.parse(ajax.responseText);
          if (data) {
            this.setState({
              currentUser: {
                uid: data.uid,
                name: data.name,
                email: data.email,
                phone: data.phone,
                roles: [...data.roles],
              },
              loggedIn: true,
            });
            window.localStorage.setItem("uid", data.uid);
            window.localStorage.setItem("name", data.name);
            window.localStorage.setItem("email", data.email);
            window.localStorage.setItem("phone", data.phone);
            window.localStorage.setItem("roles", [...data.roles]);
          }
        }
      };
    }
  }
  componentDidUpdate(_pprops, pstate) {
    if (this.state.loggedIn && !pstate.loggedIn) {
      this.updateDoctorsList();
    }
  }
  updateDoctorsList() {
    let ajax = new XMLHttpRequest();
    ajax.open("GET", "https://dev.local/api/doctors", true);
    ajax.withCredentials = true;
    ajax.setRequestHeader("Content-type", "application/json");
    ajax.send();
    ajax.onreadystatechange = () => {
      if (ajax.readyState === 4 && ajax.status === 200) {
        if (!ajax.responseText) {
          console.log("Received an empty doctors list");
          window.alert("Could not get doctors list");
        }
        let data = JSON.parse(ajax.responseText);
        if (data) {
          this.setState({
            doctors: [...data],
          });
        }
      }
      if (ajax.readyState === 4 && ajax.status !== 200) {
        console.log(ajax.responseText);
        window.alert("Could not get doctors list");
      }
    };
  }
  updateLogin(isLoggedIn) {
    this.setState({
      loggedIn: isLoggedIn,
    });
  }
  updateUser(user) {
    if (user) {
      this.setState({
        currentUser: { ...user },
      });
    } else {
      this.setState({
        currentUser: null,
      });
    }
    this.setLocalStorage(user);
  }
  setLocalStorage(user) {
    if (!user) {
      window.localStorage.removeItem("uid");
      window.localStorage.removeItem("name");
      window.localStorage.removeItem("email");
      window.localStorage.removeItem("phone");
      window.localStorage.removeItem("roles");
    } else {
      window.localStorage.setItem("uid", user.uid);
      window.localStorage.setItem("name", user.name);
      window.localStorage.setItem("email", user.email);
      window.localStorage.setItem("phone", user.phone);
      window.localStorage.setItem("roles", [...user.roles]);
    }
  }
  updateCurrentPatient(p) {
    if (!p) {
      this.setState({
        currentPatient: null,
      });
    } else {
      this.setState({
        currentPatient: { ...p },
      });
    }
  }
  logoutHandler() {
    let ajax = new XMLHttpRequest();
    let url = "https://dev.local/api/users/logout";
    ajax.open("POST", url, true);
    ajax.withCredentials = true;
    ajax.setRequestHeader("Content-type", "application/json");
    ajax.send();
    ajax.onreadystatechange = () => {
      if (ajax.readyState === 4 && ajax.status === 200) {
        // update state & localStorage
        this.setState({
          currentUser: null,
          loggedIn: false,
        });
        window.localStorage.removeItem("uid");
        window.localStorage.removeItem("name");
        window.localStorage.removeItem("email");
        window.localStorage.removeItem("phone");
        window.localStorage.removeItem("roles");
      }
    };
  }
  isAdminOrClerk() {
    if (!this.state.currentUser) {
      return false;
    } else {
      return (
        this.state.currentUser.roles.includes("Clerk") ||
        this.state.currentUser.roles.includes("Admin")
      );
    }
  }
  render() {
    return (
      <Router>
        <div className="App">
          <header>
            Sked(uler)
            <nav className="Navigation">
              <div id="logout" hidden={!this.state.loggedIn}>
                <ul>
                  <li>
                    <Link to="/">Home</Link>
                  </li>
                  <li>
                    <Link to="/Appointments">Appointments</Link>
                  </li>
                  <li>
                    <Link to="/Surgeries">Surgeries</Link>
                  </li>
                  <li>
                    <Link to="/Patients">Patients</Link>
                  </li>
                  {this.isAdminOrClerk() && (
                    <li>
                      <Link to="/Users">Users</Link>
                    </li>
                  )}
                  <li>
                    <a
                      onClick={() => {
                        this.logoutHandler();
                      }}
                    >
                      <em>Logout</em>
                    </a>
                  </li>
                </ul>
              </div>
            </nav>
          </header>
          <main>
            {!this.state.loggedIn ? (
              <Auth
                updateLogin={(s) => this.updateLogin(s)}
                updateUser={(u) => this.updateUser(u)}
                loggedIn={this.state.loggedIn}
              />
            ) : (
              <Switch>
                <Route path="/" exact>
                  <h1>Skeduler</h1>
                  <ScheduleList></ScheduleList>
                </Route>
                <Route path="/Appointments">
                  <Appointments
                    currentUser={this.state.currentUser}
                    currentPatient={this.state.currentPatient}
                    updateCurrentPatient={(p) => {
                      this.updateCurrentPatient(p);
                    }}
                    doctors={this.state.doctors}
                  />
                </Route>
                <Route path="/Surgeries">
                  <Surgeries
                    currentUser={this.state.currentUser}
                    doctors={this.state.doctors}
                  />
                </Route>
                <Route path="/Patients">
                  <Patients
                    currentPatient={this.state.currentPatient}
                    updateCurrentPatient={(p) => {
                      this.updateCurrentPatient(p);
                    }}
                  />
                </Route>
                <Route path="/Users">
                  {() => {
                    if (this.isAdminOrClerk()) {
                      return <Users />;
                    } else {
                      return null;
                    }
                  }}
                </Route>
              </Switch>
            )}
          </main>
          <footer>&#169; CG - 2020</footer>
        </div>
      </Router>
    );
  }
}
export default App;

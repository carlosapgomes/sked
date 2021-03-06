import React, { Component, Suspense, lazy } from "react";
import { BrowserRouter as Router, Switch, Route, Link } from "react-router-dom";
import "./App.css";
import Auth from "./Auth/Auth";
import { withTranslation } from "react-i18next";
const ScheduleList = lazy(() => import("./ScheduleList/ScheduleList"));
const Users = lazy(() => import("./Users/Users"));
const Patients = lazy(() => import("./Patients/Patients"));
const Surgeries = lazy(() => import("./Surgeries/Surgeries"));
const Appointments = lazy(() => import("./Appointments/Appointments"));
const ResetPassword = lazy(() => import("./ResetPassword/ResetPassword"));

// loading component for suspense fallback
class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      currentUser: null,
      loggedIn: false,
      currentPatient: null,
      currentDoctor: null,
      doctors: [],
      menuOpen: false,
    };
  }
  componentDidMount() {
    // try to recover user data in case of a page reload
    if (window.localStorage.getItem("email") !== "") {
      let ajax = new XMLHttpRequest();
      let url = "/api/users?email=" + window.localStorage.getItem("email");
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
    ajax.open("GET", "/api/doctors", true);
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
    let url = "/api/users/logout";
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
  toggleMenu() {
    this.setState({
      menuOpen: !this.state.menuOpen,
    });
  }
  render() {
    const { t } = this.props;
    return (
      <Router>
        <div className="App">
          <header>
            <nav id="logout" className="Navigation">
              <Link className="Logo" to="/">
                <img alt="Logo" src="img/sked-new.png" width="70" height="70" />
              </Link>
              <ul className={this.state.menuOpen ? "OpenedMenu" : "Menu"}>
                {this.state.loggedIn && (
                  <li>
                    <Link
                      className={
                        this.state.menuOpen ? "NavItemOpened" : "NavItem"
                      }
                      to="/Appointments"
                      onClick={() => {
                        if (this.state.menuOpen) {
                          this.toggleMenu();
                        }
                      }}
                    >
                      {t("Appointments")}
                    </Link>
                  </li>
                )}

                {this.state.loggedIn && (
                  <li>
                    <Link
                      className={
                        this.state.menuOpen ? "NavItemOpened" : "NavItem"
                      }
                      to="/Surgeries"
                      onClick={() => {
                        if (this.state.menuOpen) {
                          this.toggleMenu();
                        }
                      }}
                    >
                      {t("Surgeries")}
                    </Link>
                  </li>
                )}
                {this.state.loggedIn && (
                  <li>
                    <Link
                      className={
                        this.state.menuOpen ? "NavItemOpened" : "NavItem"
                      }
                      to="/Patients"
                      onClick={() => {
                        if (this.state.menuOpen) {
                          this.toggleMenu();
                        }
                      }}
                    >
                      {t("Patients")}
                    </Link>
                  </li>
                )}
                {this.isAdminOrClerk() && this.state.loggedIn && (
                  <li>
                    <Link
                      className={
                        this.state.menuOpen ? "NavItemOpened" : "NavItem"
                      }
                      to="/Users"
                      onClick={() => {
                        if (this.state.menuOpen) {
                          this.toggleMenu();
                        }
                      }}
                    >
                      {t("Users")}
                    </Link>
                  </li>
                )}
                {this.state.loggedIn && (
                  <li>
                    <a
                      className={
                        this.state.menuOpen ? "NavItemOpened" : "NavItem"
                      }
                      href="#!"
                      onClick={() => {
                        this.logoutHandler();
                        if (this.state.menuOpen) {
                          this.toggleMenu();
                        }
                      }}
                    >
                      <em>Logout</em>
                    </a>
                  </li>
                )}
              </ul>
              {this.state.loggedIn && !this.state.menuOpen && (
                <button
                  className="MenuToggleBtn"
                  onClick={() => {
                    this.toggleMenu();
                  }}
                >
                  Menu
                </button>
              )}
              {this.state.loggedIn && this.state.menuOpen && (
                <button
                  className="MenuToggleBtn"
                  onClick={() => {
                    if (this.state.menuOpen) {
                      this.toggleMenu();
                    }
                  }}
                >
                  X
                </button>
              )}

              {!this.state.loggedIn && (
                <a className="LoginBtn" href="/">
                  <em>Login</em>
                </a>
              )}
            </nav>
          </header>
          <main>
            {!this.state.loggedIn ? (
              <Switch>
                <Route
                  path="/ResetPassword"
                  exact
                  component={ResetPassword}
                ></Route>
                <Route path="/">
                  <Auth
                    updateLogin={(s) => this.updateLogin(s)}
                    updateUser={(u) => this.updateUser(u)}
                    loggedIn={this.state.loggedIn}
                  />{" "}
                </Route>
              </Switch>
            ) : (
              <Suspense fallback={<div>Loading...</div>}>
                <Switch>
                  <Route path="/" exact component={ScheduleList}></Route>
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
                      currentPatient={this.state.currentPatient}
                      updateCurrentPatient={(p) => {
                        this.updateCurrentPatient(p);
                      }}
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
              </Suspense>
            )}
          </main>
          <footer>
            <hr />
            <small>&#169; CG - 2020</small>
          </footer>
        </div>
      </Router>
    );
  }
}
export default withTranslation()(App);

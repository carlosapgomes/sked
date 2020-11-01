import React, { Component } from "react";
import { BrowserRouter as Router, Switch, Route, Link } from "react-router-dom";
import "./App.css";
import Auth from "./Auth/Auth";
import Appointments from "./Appointments/Appointments";
import Surgeries from "./Surgeries/Surgeries";
import Patients from "./Patients/Patients";
import Users from "./Users/Users";
class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      currentUser: null,
      loggedIn: false,
      currentPatient: null,
      currentDoctor: null,
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
  }
  logoutHandler() {
    this.setState({
      currentUser: null,
      loggedIn: false,
    });
  }
  getResult() {
    let ajax = new XMLHttpRequest();
    ajax.open("GET", "https://dev.local/api/patients", true);
    ajax.withCredentials = true;
    ajax.setRequestHeader("Content-type", "application/json");
    ajax.send();
    ajax.onreadystatechange = () => {
      if (ajax.readyState === 4 && ajax.status === 200) {
        var data = JSON.parse(ajax.responseText);
        if (typeof data != "undefined") {
          this.setState({
            result: data,
          });
        }
      }
    };
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
                  <li>
                    <Link to="/Users">Users</Link>
                  </li>
                  <li>
                    <button
                      onClick={() => {
                        this.logoutHandler();
                      }}
                    >
                      Logout
                    </button>
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
                </Route>
                <Route path="/Appointments">
                  <Appointments />
                </Route>
                <Route path="/Surgeries">
                  <Surgeries />
                </Route>
                <Route path="/Patients">
                  <Patients />
                </Route>
                <Route path="/Users">
                  <Users />
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

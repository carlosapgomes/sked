import React, { Component } from "react";
import "./App.css";
import Auth from "./Auth/Auth";
class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      currentUser: null,
      loggedIn: false,
      currentPatient: null,
      currentDoctor: null,
      result: null,
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
      <div className="App">
        {!this.state.loggedIn ? (
          <Auth
            updateLogin={(s) => this.updateLogin(s)}
            updateUser={(u) => this.updateUser(u)}
            loggedIn={this.state.loggedIn}
          />
        ) : (
          <h1>Welcome to sked {this.state.currentUser.name}</h1>
        )}
        <div>
          <button
            onClick={() => {
              this.getResult();
            }}
          >
            Query
          </button>
          <p>{this.state.result}</p>
        </div>
        <div id="logout" hidden={!this.state.loggedIn}>
          <button
            onClick={() => {
              this.logoutHandler();
            }}
          >
            Logout
          </button>
        </div>
      </div>
    );
  }
}
export default App;

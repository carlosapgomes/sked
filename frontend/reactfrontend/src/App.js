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
    };
  }
  updateLogin(isLoggedIn) {
    this.setState({
      loggedIn: isLoggedIn,
    });
  }
  updateUser(user) {
    this.setState({
      currentUser: { ...user },
    });
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
      </div>
    );
  }
}
export default App;

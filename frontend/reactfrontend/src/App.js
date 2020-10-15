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
  render() {
    return <div className="App">{!this.state.loggedIn ? <Auth /> : null}</div>;
  }
}
export default App;

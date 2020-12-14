import React from "react";
import { render } from "@testing-library/react";
import Auth from "./Auth";

it("renders without crashing", () => {
  const div = document.createElement("div");
  render(<Auth loggedIn={false} />, div);
});

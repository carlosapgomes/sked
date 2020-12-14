import React from "react";
import { render } from "@testing-library/react";
import Patients from "./Patients";

it("renders without crashing", () => {
  const div = document.createElement("div");
  render(<Patients currentPatient={null} />, div);
});

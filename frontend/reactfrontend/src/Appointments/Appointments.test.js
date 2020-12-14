import React from "react";
import { render } from "@testing-library/react";
import Appointments from "./Appointments";

it("renders without crashing", () => {
  const div = document.createElement("div");
  render(<Appointments doctors={["doc1", "doc2"]} />, div);
});

import React from "react";
import { render } from "@testing-library/react";
import Surgeries from "./Surgeries";

it("renders without crashing", () => {
  const div = document.createElement("div");
  render(<Surgeries doctors={["doc1", "doc2"]} />, div);
});

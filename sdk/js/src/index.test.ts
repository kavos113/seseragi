import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import * as fs from "fs";
import { InputData, run } from "./index";

vi.mock("fs");

describe("InputData", () => {
  it("should get the correct node data", () => {
    const data = [
      { name: "node1", data: "value1" },
      { name: "node2", data: { key: "value2" } },
    ];
    const input = new InputData(data);

    expect(input.get("node1")).toBe("value1");
    expect(input.get("node2")).toEqual({ key: "value2" });
    expect(input.get("node3")).toBeNull();
  });
});

describe("run", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("should run the handler and write output", async () => {
    const inPath = "/tmp/input.json";
    const outPath = "/tmp/output.json";
    process.env.WORKFLOW_INPUT_PATH = inPath;
    process.env.WORKFLOW_OUTPUT_PATH = outPath;

    const mockInputData = [{ name: "node1", data: "hello" }];
    vi.mocked(fs.readFileSync).mockReturnValue(JSON.stringify(mockInputData));

    const handler = async (input: InputData) => {
      const val = input.get("node1");
      return `${val} world`;
    };

    await run(handler);

    expect(fs.readFileSync).toHaveBeenCalledWith(inPath, "utf-8");
    expect(fs.writeFileSync).toHaveBeenCalledWith(outPath, JSON.stringify("hello world"));
  });

  it("should write error on failure", async () => {
    const inPath = "/tmp/input.json";
    const outPath = "/tmp/output.json";
    process.env.WORKFLOW_INPUT_PATH = inPath;
    process.env.WORKFLOW_OUTPUT_PATH = outPath;

    vi.mocked(fs.readFileSync).mockImplementation(() => {
      throw new Error("read error");
    });

    const handler = async () => "result";

    await run(handler);

    expect(fs.writeFileSync).toHaveBeenCalledWith(outPath, "Error: read error");
  });
});

import * as fs from "fs";

export type Empty = Record<string, never>;

export interface NodeData {
  name: string;
  data: any;
}

export class InputData {
  private data: NodeData[];

  constructor(data: NodeData[]) {
    this.data = data;
  }

  get<T = any>(name: string): T | null {
    const node = this.data.find((d) => d.name === name);
    if (!node) {
      return null;
    }
    return node.data as T;
  }
}

const writeError = (path: string, error: string) => {
  try {
    fs.writeFileSync(path, error);
  } catch (e) {
    console.error(`Failed to write error to ${path}:`, e);
  }
};

export async function run<O>(
  handler: (input: InputData) => Promise<O>,
): Promise<void> {
  const inPath = process.env.WORKFLOW_INPUT_PATH;
  const outPath = process.env.WORKFLOW_OUTPUT_PATH;

  if (!inPath || !outPath) {
    console.error("WORKFLOW_INPUT_PATH and WORKFLOW_OUTPUT_PATH must be set");
    return;
  }

  try {
    const inData = fs.readFileSync(inPath, "utf-8");

    const data = JSON.parse(inData) as NodeData[];
    const input = new InputData(data);

    const output = await handler(input);
    fs.writeFileSync(outPath, JSON.stringify(output));
  } catch (error) {
    console.error("Error during workflow execution:", error);
    writeError(
      outPath,
      `Error: ${error instanceof Error ? error.message : String(error)}`,
    );
  }
}

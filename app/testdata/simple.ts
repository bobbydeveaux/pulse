// A simple TypeScript file for testing

function greet(name: string): string {
  return `Hello, ${name}`;
}

function calculate(x: number, y: number, op: string): number {
  if (op === "add") {
    return x + y;
  } else if (op === "sub") {
    return x - y;
  } else if (op === "mul") {
    return x * y;
  } else if (op === "div") {
    if (y === 0) {
      throw new Error("Division by zero");
    }
    return x / y;
  }
  throw new Error(`Unknown op: ${op}`);
}

const processItems = (items: (string | number | null)[]): (string | number)[] => {
  const results: (string | number)[] = [];
  for (const item of items) {
    if (item === null) {
      continue;
    }
    if (typeof item === "string" && item.length > 0) {
      results.push(item.toUpperCase());
    } else if (typeof item === "number") {
      results.push(item * 2);
    }
  }
  return results;
};

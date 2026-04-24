# A simple Python file for testing

def hello(name):
    """Greet someone."""
    return f"Hello, {name}"


def calculate(x, y, operation):
    if operation == "add":
        return x + y
    elif operation == "sub":
        return x - y
    elif operation == "mul":
        return x * y
    elif operation == "div":
        if y == 0:
            raise ValueError("Cannot divide by zero")
        return x / y
    else:
        raise ValueError(f"Unknown operation: {operation}")


def process_items(items):
    results = []
    for item in items:
        if item is None:
            continue
        if isinstance(item, str) and len(item) > 0:
            results.append(item.upper())
        elif isinstance(item, int):
            results.append(item * 2)
    return results

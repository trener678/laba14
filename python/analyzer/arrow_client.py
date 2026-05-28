from __future__ import annotations

import argparse
import urllib.request

import pandas as pd
import pyarrow.ipc as ipc


def fetch_arrow(url: str) -> pd.DataFrame:
    with urllib.request.urlopen(url, timeout=10) as response:
        with ipc.open_stream(response.read()) as reader:
            table = reader.read_all()
    return table.to_pandas()


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--url", default="http://localhost:8080/arrow/aggregates")
    args = parser.parse_args()
    frame = fetch_arrow(args.url)
    print(frame.tail(20).to_string(index=False))


if __name__ == "__main__":
    main()

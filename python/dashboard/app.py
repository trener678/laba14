from __future__ import annotations

import time

import pandas as pd
import plotly.express as px
import streamlit as st

from analyzer.arrow_client import fetch_arrow


st.set_page_config(page_title="Laba 14 Health Dashboard", layout="wide")
st.title("Health telemetry")

endpoint = st.sidebar.text_input("Arrow endpoint", "http://localhost:8080/arrow/aggregates")
refresh = st.sidebar.slider("Refresh, seconds", min_value=1, max_value=30, value=5)

placeholder = st.empty()
while True:
    with placeholder.container():
        try:
            data = fetch_arrow(endpoint)
            if data.empty:
                st.info("Waiting for aggregates")
            else:
                data["window_start"] = pd.to_datetime(data["window_start_unix"], unit="s")
                cols = st.columns(3)
                cols[0].metric("Batches", len(data))
                cols[1].metric("Metrics", data["metric"].nunique())
                cols[2].metric("Last average", f"{data['average'].iloc[-1]:.2f}")
                st.plotly_chart(
                    px.line(data, x="window_start", y="average", color="metric"),
                    use_container_width=True,
                )
                st.dataframe(data.tail(50), use_container_width=True)
        except Exception as exc:
            st.error(str(exc))
    time.sleep(refresh)

import csv
import os
import matplotlib.pyplot as plt
import pandas as pd
import numpy as np


def plotDrops(fig: plt.Figure, ax: plt.Axes, df: pd.DataFrame, agent_name: str, color: str, dir: str):
    ax.plot(df["cbr_rate"], df["avg_drops"], marker='o', label=agent_name, color=color)
    # Shade in standard deviation
    ax.fill_between(df["cbr_rate"], df["avg_drops"] - df["std_drops"], df["avg_drops"] + df["std_drops"], alpha=0.25, color=color)
    ax.legend()
    ax.set_title("TCP Dropped Packets vs. CBR Rate")
    ax.set_xlabel('CBR Rate (Mbps)')
    ax.set_ylabel('TCP Dropped Packets')
    fig.savefig("{}/exp01_drops.png".format(dir))

def plotLatency(fig: plt.Figure, ax: plt.Axes, df: pd.DataFrame, agent_name: str, color: str, dir: str):
    ax.plot(df["cbr_rate"], df["avg_latency"]*1000, marker='o', label=agent_name, color=color)
    # Shade in standard deviation
    ax.fill_between(df["cbr_rate"], df["avg_latency"]*1000 - df["std_latency"]*1000, df["avg_latency"]*1000 + df["std_latency"]*1000, alpha=0.25, color=color)
    ax.legend()
    ax.set_title("TCP Latency vs. CBR Rate")
    ax.set_xlabel('CBR Rate (Mbps)')
    ax.set_ylabel('TCP Latency (ms)')
    fig.savefig("{}/exp01_latency.png".format(dir))

def plotThroughput(fig: plt.Figure, ax: plt.Axes, df: pd.DataFrame, agent_name: str, color: str, dir: str):
    ax.plot(df["cbr_rate"], df["avg_throughput"], marker='o', label=agent_name, color=color)
    # Shade in standard deviation
    ax.fill_between(df["cbr_rate"], df["avg_throughput"] - df["std_throughput"], df["avg_throughput"] + df["std_throughput"], alpha=0.25, color=color)
    ax.legend()
    ax.set_title("TCP Throughput vs. CBR Rate")
    ax.set_xlabel('CBR Rate (Mbps)')
    ax.set_ylabel('TCP Throughput (Mbps)')
    ax.set_yticks(np.arange(0, 10, 1))
    fig.savefig("{}/exp01_throughput.png".format(dir))
    

def main():

    plt.style.use('ggplot')

    # Init the plots outside of the loop so we can share them
    # Subplots are stateful, so it's kind of like a pointer
    fig1, ax1 = plt.subplots()
    fig2, ax2 = plt.subplots()
    fig3, ax3 = plt.subplots()

    directory = "../results/exp01"
    csvfiles = ["exp01_Tahoe.csv", "exp01_Reno.csv", "exp01_Newreno.csv", "exp01_Vegas.csv"]
    
    for csvfile in csvfiles:
        # The agent_name is between the '_' and the '.csv'
        agent_name = csvfile.split('_')[1].split('.')[0]
        colorMap = {"Tahoe": "tab:red",
                    "Reno": "tab:orange",
                    "Newreno": "tab:green",
                    "Vegas": "tab:blue"}
        
        with open(os.path.join(directory, csvfile), 'r') as f:
            reader = csv.reader(f)
            # Grab the header row
            header = next(reader)
            data = list(reader)
            df = pd.DataFrame(data, columns=header)
            for col in df.columns:
                df[col] = df[col].apply(lambda x: round(float(x), 10)) # Round all values to 10 decimal places

            plotDrops(fig1, ax1, df, agent_name, colorMap[agent_name], directory)
            plotLatency(fig2, ax2, df, agent_name, colorMap[agent_name], directory)
            plotThroughput(fig3, ax3, df, agent_name, colorMap[agent_name], directory)

if __name__ == "__main__":
    main()
import csv
import os
import matplotlib.pyplot as plt
import pandas as pd
import numpy as np


def plotDrops(df: pd.DataFrame, agent1: str, agent2: str, color1: str, color2: str, save_dir: str):
    fig, ax = plt.subplots()
    ax.plot(df["cbr_rate"], df["avg_drops1"], marker='o', label=agent1, color=color1)
    ax.plot(df["cbr_rate"], df["avg_drops2"], marker='o', label=agent2, color=color2)
    # Shade in standard deviation
    ax.fill_between(df["cbr_rate"], df["avg_drops1"] - df["std_drops1"], df["avg_drops1"] + df["std_drops1"], alpha=0.25, color=color1)
    ax.fill_between(df["cbr_rate"], df["avg_drops2"] - df["std_drops2"], df["avg_drops2"] + df["std_drops2"], alpha=0.25, color=color2)
    # Set the legend to the left side
    ax.legend(loc='upper left')
    ax.set_title(agent1 + "/" + agent2 + " Dropped Packets vs. CBR Rate")
    ax.set_xlabel('CBR Rate (Mbps)')
    ax.set_ylabel('TCP Dropped Packets')
    ax.set_yticks(np.arange(0, 130, 20))
    savefile = "{}/exp02_{}_{}_drops.png".format(save_dir, agent1, agent2)
    fig.savefig(savefile)

def plotLatency(df: pd.DataFrame, agent1: str, agent2: str, color1: str, color2: str, save_dir: str):
    fig, ax = plt.subplots()
    ax.plot(df["cbr_rate"], df["avg_latency1"]*1000, marker='o', label=agent1, color=color1)
    ax.plot(df["cbr_rate"], df["avg_latency2"]*1000, marker='o', label=agent2, color=color2)
    # Shade in standard deviation
    ax.fill_between(df["cbr_rate"], df["avg_latency1"]*1000 - df["std_latency1"]*1000, df["avg_latency1"]*1000 + df["std_latency1"]*1000, alpha=0.25, color=color1)
    ax.fill_between(df["cbr_rate"], df["avg_latency2"]*1000 - df["std_latency2"]*1000, df["avg_latency2"]*1000 + df["std_latency2"]*1000, alpha=0.25, color=color2)
    ax.legend()
    ax.set_title(agent1 + "/" + agent2 + " Latency vs. CBR Rate")
    ax.set_xlabel('CBR Rate (Mbps)')
    ax.set_ylabel('TCP Latency (ms)')
    # ax.set_yticks(np.arange(30, 55, 5))
    savefile = "{}/exp02_{}_{}_latency.png".format(save_dir, agent1, agent2)
    fig.savefig(savefile)

def plotThroughput(df: pd.DataFrame, agent1: str, agent2: str, color1: str, color2: str, save_dir: str):
    fig, ax = plt.subplots()
    ax.plot(df["cbr_rate"], df["avg_throughput1"], marker='o', label=agent1, color=color1)
    ax.plot(df["cbr_rate"], df["avg_throughput2"], marker='o', label=agent2, color=color2)
    # Shade in standard deviation
    ax.fill_between(df["cbr_rate"], df["avg_throughput1"] - df["std_throughput1"], df["avg_throughput1"] + df["std_throughput1"], alpha=0.25, color=color1)
    ax.fill_between(df["cbr_rate"], df["avg_throughput2"] - df["std_throughput2"], df["avg_throughput2"] + df["std_throughput2"], alpha=0.25, color=color2)
    ax.legend()
    ax.set_title(agent1 + "/" + agent2 + " Throughput vs. CBR Rate")
    ax.set_xlabel('CBR Rate (Mbps)')
    ax.set_ylabel('TCP Throughput (Mbps)')
    ax.set_yticks(np.arange(0, 10, 1))
    savefile = "{}/exp02_{}_{}_throughput.png".format(save_dir, agent1, agent2)
    fig.savefig(savefile)
    

def main():

    plt.style.use('ggplot')

    dir = "../results/exp02"
    csvfiles = ["exp02_Reno_Reno.csv", "exp02_Newreno_Reno.csv", "exp02_Newreno_Vegas.csv", "exp02_Vegas_Vegas.csv"]

    colorMap = {"Tahoe": "tab:red",
                "Reno": "tab:orange",
                "Newreno": "tab:green",
                "Vegas": "tab:blue"}
    
    for csvfile in csvfiles:
        # The agent_name is between the first '_' and the '.csv'
        agent_pair = csvfile.split('_', 1)[1].split('.')[0]
        agent1, agent2 = agent_pair.split('_')
        
        with open(os.path.join(dir, csvfile), 'r') as f:
            reader = csv.reader(f)
            # Grab the header row
            header = next(reader)
            data = list(reader)
            df = pd.DataFrame(data, columns=header)
            for col in df.columns:
                df[col] = df[col].apply(lambda x: round(float(x), 10)) # Round all values to 10 decimal places

            color1 = colorMap[agent1]
            if agent1 == agent2:
                if agent2 == "Reno":
                    color2 = "tab:purple" # orange and purple
                else: # Vegas
                    color2 = "tab:red" # blue and red
            else:
                color2 = colorMap[agent2]

            plotDrops(df, agent1, agent2, color1, color2, dir)
            plotLatency(df, agent1, agent2, color1, color2, dir)
            plotThroughput(df, agent1, agent2, color1, color2, dir)

if __name__ == "__main__":
    main()
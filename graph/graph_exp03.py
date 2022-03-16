import csv
import os
import matplotlib.pyplot as plt
import pandas as pd

def main():

    plt.style.use('ggplot')

    dir = "../results/exp03"

    groups = [
        ["exp03_Reno_DropTail_CBR.csv", "exp03_Reno_DropTail_TCP.csv"],
        ["exp03_Reno_RED_CBR.csv", "exp03_Reno_RED_TCP.csv"],
        ["exp03_Sack1_DropTail_CBR.csv", "exp03_Sack1_DropTail_TCP.csv"],
        ["exp03_Sack1_RED_CBR.csv", "exp03_Sack1_RED_TCP.csv"]
    ]

    for group in groups:
        fig, ax = plt.subplots()
        for csvfile in group:

            fname = csvfile.split('.')[0].split('_')
            agent = fname[1]
            queue = fname[2]
            flow = fname[3]
            
            with open(os.path.join(dir, csvfile), 'r') as f:
                reader = csv.reader(f)
                # Grab the header row
                header = next(reader)
                data = list(reader)
                df = pd.DataFrame(data, columns=header)

                for col in df.columns:
                    df[col] = df[col].apply(lambda x: round(float(x), 3)) # Round all values to 3 decimal places

                if flow == "CBR":
                    ax.plot(df['time_ticks'], df['throughput_ticks'], color='tab:blue', label='CBR')
                else:
                    ax.plot(df['time_ticks'], df['throughput_ticks'], color='tab:red', label='TCP')
            
            ax.legend()
            ax.set_xlabel('Time (s)')
            ax.set_ylabel('Throughput (Mbps)')
            ax.set_title("{}/{} Throughput over Time".format(agent, queue))
            filename = "{}/{}_{}_trace.png".format(dir, agent, queue)
            fig.savefig(filename)

            

if __name__ == "__main__":
    main()
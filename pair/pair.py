import pandas as pd
import matplotlib.pyplot as plt


class Pair:
    def __init__(self, pairs: pd.DataFrame, threshold: float = None, fee: float = None):
        self.pairs = pairs
        self.threshold = threshold
        self.fee = fee

    def run(self):
        if self.threshold is None:
            return self._none_threshold_run()
        return None

    def _none_threshold_run(self):
        columns = self.pairs.columns
        diff: pd.Series = self.pairs[columns[0]] - self.pairs[columns[1]]
        threshold = diff.abs().mean()
        plt.plot(diff.index, diff)
        plt.plot(diff.index, [threshold] * len(diff))
        sign = diff.apply(lambda x: 1 if x > 0 else -1)
        trading_flags: pd.Series = diff.apply(lambda x: abs(x) // threshold) * sign
        trading_action = trading_flags.diff(1)
        trading_action[0] = diff[0] // threshold
        profit = (trading_action * diff).sum()
        plt.show()
        return profit

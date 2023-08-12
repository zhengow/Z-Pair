import os

from pair import Pair
import pandas as pd
import numpy as np
import sqlite3
import dolphindb as ddb
from concurrent.futures import ThreadPoolExecutor


def read_sql(coin: str):
    return '''
    SELECT symbol as instrument, datetime, volume, open_price as open, close_price as close, high_price as high, low_price as low 
    FROM dbbardata
    where symbol='{}' and interval='1m' and exchange='BINANCE'
    '''.format(coin)


def mock():
    conn = sqlite3.connect("D:\\work\\vnpy-work\\.vntrader\\database.db")
    s = ddb.session()
    s.connect("localhost", 8848, "admin", "123456")
    # df = pd.read_sql_query("select distinct symbol from dbbardata where interval='1m' and exchange = 'BINANCE'", conn)
    df = pd.read_pickle('symbol.pkl')
    print(df['symbol'])
    df.to_pickle('symbol.pkl')
    # def insert(coin):
    #     print(f"start insert {coin}")
    #     conn2 = sqlite3.connect("D:\\work\\vnpy-work\\.vntrader\\database.db")
    #     sql = read_sql(coin)
    #     df = pd.read_sql_query(sql, conn2)
    #     appender = ddb.tableAppender(dbPath="dfs://bar", tableName="min", ddbSession=s)
    #     df['datetime'] = pd.to_datetime(df['datetime'])
    #     num = appender.append(df)
    #     conn2.close()
    #     return f'{coin}: {num}'
    # with ThreadPoolExecutor(max_workers=10) as pool:
    #     results = pool.map(insert, df['symbol'])
    # for r in results:
    #     print(r)
    for idx, symbol in enumerate(df['symbol']):
        print(f"start read {symbol}")
        f = f'{symbol}.pkl'
        if os.path.exists(f):
            df2 = pd.read_pickle(f)
        else:
            sql = read_sql(symbol)
            df2 = pd.read_sql_query(sql, conn)
            df2.to_pickle(f)
        print(df2)
        print(f"finish read {symbol}")
        df2['datetime'] = pd.to_datetime(df2['datetime'])
        # cwd = os.getcwd()
        # f = os.path.join(cwd, f"{symbol}.pkl")
        # s.loadTextEx(dbPath="dfs://bar", tableName='min', partitionColumns=["instrument", "datetime"], remoteFilePath=f)
        appender = ddb.tableAppender(dbPath="dfs://bar", tableName="min", ddbSession=s)
        num = appender.append(df2)
        print(f"finish {idx+1}/{len(df['symbol'])}, {symbol}")

    s.close()
    conn.close()

    # dates = pd.date_range(start='2023-01-01', end='2023-08-11', freq='D')
    #
    # # 创建股票价格数据
    # stock1_prices = np.random.uniform(100, 200, len(dates))
    # stock2_prices = np.random.uniform(50, 150, len(dates))
    #
    # # 创建DataFrame
    # df = pd.DataFrame({'Stock1': stock1_prices, 'Stock2': stock2_prices}, index=dates)
    return df


def main():
    df = mock()
    # p = Pair(df)
    # profit = p.run()
    # print(profit)


if __name__ == '__main__':
    main()


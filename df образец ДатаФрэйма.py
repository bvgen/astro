import pandas as pd
import numpy as np
df = pd.DataFrame({'A': 'foo bar foo bar foo bar foo foo'.split(),
                   'B': 'one one two three two two one three'.split(),
                   'C': np.arange(8), 'D': np.arange(8) * 2})
print(df)



time = ['2:00', '5:09']
t = pd.DataFrame({'Time': time})

for i in range(10): 
    t = t.append({'Time': i}, ignore_index = True)



plnt = ('Sun', 'Moon', 'Merc', 'Ven', 'Mars', 'Jup', 'Sat',
        'Uran', 'Nep', 'Plut', 'Uzel')
f1 = pd.DataFrame(columns = plnt, index = plnt)
f1 = f1.fillna(0)
print(f1)

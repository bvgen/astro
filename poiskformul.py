import pandas as pn
import math
import time as timer
start = timer.time() #запускаем секундомер
pn.set_option('display.precision', 2) #отображает 2 знака после запятой
pn.set_option('display.width', 150) #чтобы ДатаФрэйм помещался на экране
pn.set_option('display.max_columns', None) 
pn.set_option('display.max_rows', 200)
pn.options.display.float_format ='{:,.2f}'.format
#создаем таблицу с положениями планет в зависимости от времени
"""
Эфемериды выгружены из программы Zet8 или Zet9 в текстовый файл, затем в Ef5let.csv
"""
tab = pn.read_csv('Ef5let.csv', sep = ';', skiprows = range(0,2)) #
tab['Дата, время'] = tab['Unnamed: 0'] + " " + tab['Дата,'] #объединяем столбцы
del tab['Unnamed: 0']
del tab['Дата,']
new = tab['Дата, время']  #создаем Series из столбца, который надо передвинуть
                          #название должно отличаться
tab.insert(0,'Дата и время',new) #вставляем этот Series как новый столбец на позицию 0
del tab['Дата, время']           #удаляем скопированный столбец
nazv2 = ['DayTime', 'StarTime', 'Sun', 'Moon', 'Merc', 'Ven', 'Mars', 'Jup', 'Sat',
        'Uran', 'Nep', 'Plut', 'Uzel']
tab.columns = nazv2
#дату и время переводим из строк в значения времени
tab['DayTime'] = pn.to_datetime(tab['DayTime'], dayfirst = True, errors = 'coerce')
bw = pn.read_csv('bazawin.csv', sep=',')#- - - загружаем базу данных удачных случаев - - -
del bw['Unnamed: 0']
bw['DayTime'] = pn.to_datetime(bw['DayTime'], dayfirst = True, errors = 'coerce')
# делаем массив с координатами планет на момент рождения человека
radix = pn.DataFrame({'Nazv': ['r.Sun', 'r.Moon', 'r.Merc', 'r.Ven', 'r.Mars', 'r.Jup',
                               'r.Sat', 'r.Uran', 'r.Nep', 'r.Plut', 'r.Uzel',
                               'r.Asc', 'r.MC'],
                       'Koords': [96.8, 212.735, 99.8, 63.445, 132.192, 347.743, 98.12,
                                  203.663, 247.483, 184.121, 259.358, 73.54, 316.85],
                       'Speed': [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]})

def interpol_time(i1, i2, time):    
    t1 = tab.loc[i1, 'DayTime'] # 
    t2 = tab.loc[i2, 'DayTime']    
    dt = pn.Timedelta.total_seconds(t2 - t1)    
    tx = pn.Timedelta.total_seconds(time - t1)
    if dt == 0:
        k_in = 0
    else:
        k_in = tx / dt
    return k_in# - - - - коэфф линейной интерполяции - - - 

def interpol_zv_time(i1, i2, k_i):
    st1 = tab.StarTime[i1]  # в часах
    st2 = tab.StarTime[i2]
    ds = st2 - st1
    if ds < 0:
        ds = ds + 24
    dx = ds * k_i
    startime = st1 + dx #звездное время в секундах от начала суток
    if startime > 24:
        startime = startime - 24    
    return startime

def raschet_kuspidov(startime):
    stm = 360 * startime / 24    
    ramc = math.radians(stm) # расчет МС
    e = math.radians(23.44166)
    mc = math.atan(math.tan(ramc) / math.cos(e))
    k10 = math.degrees(mc)
    #проверка и выравнивание МС
    if abs(stm - k10) > 10:
        k10 = k10 + 180
    if abs(stm - k10) > 10:
        k10 = k10 + 180
    mc = math.radians(k10)
    # расчет куспидов домов по системе Коха
    lat = math.radians(43.25) #географ. широта Алматы
    dec = math.asin(math.sin(mc) * math.sin(e)) #склонение МС
    oamc1 = math.tan(dec) * math.tan(lat)  #наклонное восхождение МС
    oamcr = ramc - math.asin(oamc1) #в радианах
    oamc = math.degrees(oamcr)  #в градусах
    dx = (stm + 90 - oamc) / 3 #интервал между куспидами
    if dx < 0:
        dx = dx + 360
    h11 = oamc + dx - 90  #в градусах
    h12 = h11 + dx
    h1 = h12 + dx
    h2 = h1 + dx
    h3 = h2 + dx
    h11r = math.radians(h11) #в радианах
    h12r = math.radians(h12)
    h1r = math.radians(h1)
    h2r = math.radians(h2)
    h3r = math.radians(h3)

    #далее положение куспидов на эклиптике
    ts = math.tan(lat) * math.sin(e)
    e11r = math.atan((- ts - math.sin(h11r) * math.cos(e)) / math.cos(h11r))
    e11 = 90 - math.degrees(e11r)
    e12r = math.atan((- ts - math.sin(h12r) * math.cos(e)) / math.cos(h12r))
    e12 = 90 - math.degrees(e12r)
    e1r = math.atan((- ts - math.sin(h1r) * math.cos(e)) / math.cos(h1r))
    e1 = 90 - math.degrees(e1r)
    e2r = math.atan((- ts - math.sin(h2r) * math.cos(e)) / math.cos(h2r))
    e2 = 90 - math.degrees(e2r)
    e3r = math.atan((- ts - math.sin(h3r) * math.cos(e)) / math.cos(h3r))
    e3 = 90 - math.degrees(e3r)
    return (e1, e2, e3, 0, 0, 0, 0, 0, 0, k10, e11, e12)

def check_kusp(ku1, ku2): #вторая функция для правильной расстановки куспидов
    razn = ku2 - ku1
    if razn < -180: ku2 = ku2  
    else:
        if ku2 > ku1: ku2 = ku2
        else: ku2 = ku2 + 180
    return ku2

def align_kusp(ku): #первая функция для правильной расстановки куспидов
    n = 9
    n6 = 0
    n2 = 9
    while n != 4:
        n2 = n + 1
        if n2 == 12:
            n2 = 0
        ku[n2] = check_kusp(ku[n], ku[n2])
        n6 = n + 6
        if n6 > 11:
            n6 = n6 - 12
        ku[n6] = ku[n] + 180
        if ku[n6] > 360:
            ku[n6] = ku[n6] - 360
        n = n + 1
        if n == 12:
            n = 0
    return ku

def tab_planets(nazv2, i1, i2, k_i):
    dtm = 1 #шаг времени в DataFrame tab, dtm = 1 час
    nazv3 = nazv2[2:]
    kor = [] #список координат планет
    sk = []  #список скоростей планет, угловых секунд в час
    i = 0
    while i < 11:
        colname = nazv3[i]
        p = getkoord(i1, i2, k_i, colname, dtm)
        kx, v = p
        kor.append(kx)
        sk.append(v)  
        i = i + 1
    hoz = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, '-']
    planets = pn.DataFrame({'Nazv': nazv3, 'Koords': kor, 'Uprav': hoz, 'Speed': sk})
    #print(planets)
    return planets

def skorost_uglov(ku, startime): # аргумент - список куспидов
    dzv = 0.1 #0.1 часа (6 минут) Звездного Времени
    startime2 = startime + dzv
    asc1 = ku[0]
    mc1 = ku[9]
    kus2 = raschet_kuspidov(startime2) #расчет куспидов
    ku2 = list(kus2)
    asc2 = check_kusp(ku2[9], ku2[0])  #выравнивание Асц2
    mc2 = ku2[9]
    #считаем скорость в угловых секундах в час
    xmc = mc2 - mc1 #пробег МС в градусах
    if xmc < 0:
        xmc = xmc + 360
    vmc = 3600 * xmc / (dzv * 360 / 361)
    # 360/361 - это соотношение между локальным временем и звездным
    xasc = asc2 - asc1
    if xasc < 0:
        xasc = xasc + 360
    vasc = 3600 * xasc / (dzv * 360 / 361)
    return vasc, vmc

def skorost_uglov2(ku, shag, asc0, mc0): # аргумент - список куспидов
    asc1 = ku[0]
    mc1 = ku[9]
    xmc = mc1 - mc0 #пробег МС в градусах
    if xmc < 0:
        xmc = xmc + 360
    vmc = 36000 * 361 * xmc / shag
    # 360/361 - это соотношение между локальным временем и звездным
    xasc = asc1 - asc0
    if xasc < 0:
        xasc = xasc + 360
    vasc = 36000 * 361 * xasc / shag
    return vasc, vmc

def getkoord(i1, i2, k_i, stolbec, dtm):#получаем координаты планет из tab
    k1 = tab.loc[i1, stolbec]
    k2 = tab.loc[i2, stolbec]
    k3 = tab.loc[(i1 + 1), stolbec]    
    ds = k2 - k1 #разница в координатах
    ds2 = k3 - k1
    if ds < -300:  #в случае пересечения О град Овна
        ds = ds + 360
        ds2 = ds2 + 360
    v = 3600 * ds2 / dtm  #скорость планеты, угловых секунд в час 
    if ds < 0:
        kx = k2 - ds * k_i #применяем интерполяцию
    else:
        kx = k1 + ds * k_i
    if kx >= 360:
        kx = kx - 360
    return kx, v

def raschet_tochek(asc, sun, moon):
    den = 0
    das = asc - sun   #определение дня или ночи
    if das < -180:
        das = 360 + das
    if das > 0:
        den = 1    #это день
    else:
        den = -1   #это ночь
    fortune = asc + den * (moon - sun)
    tduha = asc + den * (sun - moon)
    if fortune > 360:
        fortune = fortune - 360
    if fortune < 0:
        fortune = fortune + 360
    if tduha > 360:
        tduha = tduha - 360
    if tduha < 0:
        tduha = tduha + 360   
    return fortune, tduha

def upr_dom(planets, ku):
    hoz_znak = ('Mars', 'Ven', 'Merc', 'Moon', 'Sun', 'Merc', 'Ven', 'Mars', 'Jup',
            'Sat', 'Sat', 'Jup')
    hoz_znak2 = ('-', '-', '-', '-', '-', '-', '-', 'Plut', '-', '-', 'Uran', 'Nep')
    dom = (1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12)
    upr = pn.DataFrame({'Dom': dom, 'Koord': ku})
    upr['Uprav1'] = upr['Koord']
    upr['Uprav2'] = dom
    n = 0
    aa = pn.Index(planets['Nazv'])
    while n < 12:
        x = int(ku[n] // 30)  #номер знака, в котором куспид
        upr.loc[n, 'Uprav1'] = hoz_znak[x]
        upr.loc[n, 'Dom'] = n + 1
        upr.loc[n, 'Uprav2'] = hoz_znak2[x]
        n = n + 1
    c = 0
    zc = []
    while c < 10:
        n = 0
        zn = []
        while n < 12:
            if planets.Nazv[c] == upr.Uprav1[n]:
                nom = upr.Dom[n]  
                zn.append(nom)
            if planets.Nazv[c] == upr.Uprav2[n]:
                nom = upr.Dom[n]  
                zn.append(nom)
            n = n + 1
        c = c + 1
        zc.append(zn)
    zc = zc + ['']
    planets['Uprav'] = zc
    return upr

def klassif_asp(ugol):  #классификация угла, к какому виду аспектов он относится
    major = [0, 60, 90, 120, 180]
    orbmain = 11#было 9
    minor = [30, 45, 135, 150]
    orbminor = 1
    #micro = [15, 75, 105, 165]
    #orbmicro = 0.25
    qu = -1
    for c in major:
        k = math.isclose(ugol, c, abs_tol = orbmain)
        if k == True:
            qu = c
            return qu
    """
    for c in minor:
        k = math.isclose(ugol, c, abs_tol = orbminor)
        if k == True:
            qu = c
            return qu
   
    for c in micro:
        k = math.isclose(ugol, c, abs_tol = orbmicro)
        if k == True:
            qu = c
            return qu
    """
    return qu

def check_shod(p1, p2, v1, v2, u, vid): #проверка аппликационности аспектов
    pr = math.radians(p1 - p2)
    pp = math.sin(pr)
    if vid == -1:
        return False
    if v1 > v2:
        if pp < 0:
            if u < vid:
                applic = False
            else:
                applic = True
        else:
            if u < vid:
                applic = True
            else:
                applic = False       
    else:
        if pp < 0:
            if u < vid:
                applic = True
            else:
                applic = False
        else:
            if u < vid:
                applic = False
            else:
                applic = True
    if applic == True:
        app = 'sho'
        return app
    app = 'ras'
    return app

def raschet_asp(planets):
    #- - - - - делаем DataFrame для аспектов - - - - - - - -
    df = pn.DataFrame(columns = ('plnt1', 'upr1', 'ugol', 'plnt2', 'upr2', 'orb', 'vid', 'appl'), )
    krd = 'Koords'
    im = 'Nazv'
    hoz = 'Uprav'
    dlin = len(planets[krd])
    sk = 'Speed'
    df['upr1'] = df['upr1'].astype('object') #чтобы мог принять список
    df['upr2'] = df['upr2'].astype('object')
    df['vid'] = df['vid'].astype('object')    
    i = 0
    z = 0
    while i < dlin: #цикл и вложенный цикл для расчета аспектов 
        p1 = planets.loc[i, krd]
        upr1 = planets.loc[i, hoz]
        im1 = planets.loc[i, im]
        v1 = planets.loc[i, sk]
        j = 0
        while j < dlin:
            if i == j:
                j = j + 1
                continue
            p2 = planets.loc[j, krd]
            u = abs(p1 - p2)
            if u > 180:
                u = 360 - u
            vidasp = klassif_asp(u)
            if vidasp == -1:
                j = j + 1
                continue
            im2 = planets.loc[j, im]
            upr2 = planets.loc[j, hoz]
            v2 = planets.loc[j, sk]
            shod = check_shod(p1, p2, v1, v2, u, vidasp)
            if shod == "ras":  #расходящиеся (сепарационные аспекты не будем записывать)
                j = j + 1
                continue
            # - - - - запись данных в DataFrame  asp - - - - - - 
            df.loc[z, 'plnt1'] = im1            
            df.at[z, 'upr1'] = upr1#для присвоения списка надо использовать at
            df.loc[z, 'ugol'] = u
            df.loc[z, 'plnt2'] = im2            
            df.at[z, 'upr2'] = upr2
            df.loc[z, 'vid'] = vidasp  
            df.loc[z, 'appl'] = shod
            df.loc[z, 'orb'] = abs(u - vidasp)
            j = j + 1
            z = z + 1
        i = i + 1    
    t2 = df.reset_index(drop = True)
    return t2
        
def add_tochki(planets, ku, vasc, vmc):
    #- - - - - рассчитываем жребии Фортуны и Духа - - - -
    sun = planets.loc[0, 'Koords']
    moon = planets.loc[1, 'Koords']
    tochki = raschet_tochek(ku[0], sun, moon)
    fortune, tduha = tochki
    #- - - - - дополняем таблицу planets - - - - - - -
    planets.loc[11] = ['Asc', ku[0], '-', vasc]
    planets.loc[12] = ['MC', ku[9], '-', vmc]
    planets.loc[13] = ['Fortune', fortune, '-', vasc]
    planets.loc[14] = ['tDuha', tduha, '-', vasc]    

def raschet_tranz(planets, radix):    
    #- - - - - расчет аспектов транзитных планет к радикальным - - - - - - - -
    tran = pn.DataFrame(columns = ('plnt1', 'upr1', 'ugol', 'plnt2', 'upr2', 'orb', 'vid', 'appl'), )
    tran['upr1'] = tran['upr1'].astype('object')
    tran['upr2'] = tran['upr2'].astype('object')
    tran['vid'] = tran['vid'].astype('object')    
    krd = 'Koords'
    im = 'Nazv'
    hoz = 'Uprav'    
    sk = 'Speed'
    i = 0
    z = 0
    while i < 15:
        t1 = planets.loc[i, krd]        
        im1 = planets.loc[i, im]
        v1 = planets.loc[i, sk]
        upr1 = planets.loc[i, hoz]
        j = 0
        while j < 13:
            t2 = radix.loc[j, krd]
            u = abs(t1 - t2)
            if u > 180:
                u = 360 - u
            vidasp = klassif_asp(u)
            if vidasp == -1:
                j = j + 1
                continue
            #if abs(u - vidasp) > 2: #для ограничения по орбису
                #j = j + 1
                #continue
            v2 = radix.loc[j, sk]
            shod = check_shod(t1, t2, v1, v2, u, vidasp)
            if shod == "ras":  #расходящиеся (сепарационные аспекты не будем записывать)
                j = j + 1
                continue
            im2 = radix.loc[j, im]            
            tran.loc[z, 'plnt1'] = im1
            tran.at[z, 'upr1'] = upr1
            tran.loc[z, 'ugol'] = u
            tran.loc[z, 'plnt2'] = im2            
            tran.loc[z, 'vid'] = vidasp  
            tran.loc[z, 'appl'] = shod
            tran.loc[z, 'orb'] = abs(u - vidasp)            
            j = j + 1
            z = z + 1
        i = i + 1
    t2 = tran.reset_index(drop = True)
    return t2


def main_calc(time, tab, radix):
    aa = pn.Index(tab['DayTime'])
    try:
        i1 = aa.get_loc(time, method = 'pad') #левая граница
        i2 = aa.get_loc(time, method = 'backfill') #правая граница
    except Exception:
        print("нет эфемерид")
        return
    if i1 != i2:
        k_i = interpol_time(i1, i2, time)
    else:
        k_i = 0    
    #- - - - - - - - - - - - - - - - - - - -
    startime = interpol_zv_time(i1, i2, k_i)    
    kus = raschet_kuspidov(startime) #запускаем расчет куспидов
    ku = list(kus)                   #кортеж куспидов преобразуем в список
    kus3 = align_kusp(ku) #расчет недостающих и выравнивание всех куспидов
    ku = list(kus3)
    vu = skorost_uglov(ku, startime)      
    vasc, vmc = vu    
    global planets, asp, asp2, tranz
    planets = tab_planets(nazv2, i1, i2, k_i) #делаем DataFrame свойств планет 
    hoz_dom = upr_dom(planets, ku)
    add_tochki(planets, ku, vasc, vmc) #добавляем точки в таблицу planets
    asp = raschet_asp(planets) #расчет аспектов и внесение их в ДатаФрэйм - -
    tranz = raschet_tranz(planets, radix)
    asp2 = asp.append(tranz)
    asp2 = asp2.reset_index(drop = True)
    return 
#- - - - - - - - - - - - - - - - - - - - - - - - -
rez = pn.DataFrame(columns = ['Sun', 'oon', 'erc', 'Ven', 'ars', 'Jup', 'Sat',
        'ran', 'Nep', 'lut', 'zel'], index = [0])
rez = rez.fillna(0)

def merc_planeta_i_ug_ne():
    global time, db, m, rez
    s = pn.DataFrame(columns = ['plnt1', 'vid', 'plnt2', 'orb'], index = [])
    ug = ('Asc', 'MC')
    bt = ('Asc', 'MC', "Fortune", 'tDuha', 'r.Asc', 'r.MC')
    u = asp2.query('plnt1 in @ug and plnt2 not in @bt')
    u = u.sort_values(by = 'orb', ascending = True, ignore_index = True)
    up = u.plnt2[0]
    if up[-3:] != "Nep":
        return
    pp = asp2.query('plnt1 not in @bt and plnt2 not in @bt')
    pp = pp.sort_values(by = 'orb', ascending = True, ignore_index = True)
    p1 = pp.plnt1[0]
    p2 = pp.plnt2[0]
    if p1[-3:] != "erc" and p2[-3:] != "erc":
        return
    s = s.append(u.loc[0, ])
    s = s.append(pp.loc[0, ])
    if p1[-3:] == "erc":
        px = p2[-3:]
    else:
        px = p1[-3:]
    rez.loc[0, px] = rez.loc[0, px] + 1    
    print('\nВремя в Алмате ', time, "Удачное событие")
    print(xt, "из ", db)
    print(s)
    m = m + 1
    return            
    
#- - - - - - - - - - - - - - - - - - - - - - - - -
#- - - - - - - - - - - - - - - - - - - - - - - - -
baz = bw #bw - это таблица событий, с датой и описанием события
db = len(baz)
xt, m = 0, 0
while xt < db:
    time = baz.DayTime[xt]
    #в оригинале des = baz.Description[xt]    
    main_calc(time, tab, radix)
    # - - - - - - - - - - - - - - - - -
    print(asp.loc[:19, ])
    merc_planeta_i_ug_ne()     #сочетание планет, которое ищем
    # - - - - - - - - - - - - - - - - - 
    xt = xt + 1
    #print(xt, "из ", db)
    
# - - - - конец цикла - - - - - -
print('Найдено = ', m)
print(rez) #вывод результатов подсчета
finish = timer.time()
print("Продолжительность", finish - start, "с")
input("нажмите Enter")

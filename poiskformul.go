package main
//- - - - - - расширенная функция Ме-Х + Угол-Нептун - - - - - -
//Сочетание: самый точный угол в аспекте с Нептуном, 
//а самый точный межпланетный аспект - это Меркурий с другой планетой
//подсчитываем количество встреченных других планет 
import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"sort"
	"math"
	"reflect"
)
func main() {
	start := time.Now()   // запускаем секундомер
	Makeefm()	
	cycle()
	fin(start)
}
func fin(start time.Time){
	t2 := time.Now()
	elapsed := t2.Sub(start)
	fmt.Println(elapsed)
	var bukva string
	fmt.Println("press Enter")
	fmt.Scanf("%s\n", &bukva)
}
var summ[11] int
var names[11] string
var itogo int
func cycle() {	
	summ = [11]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
        names = [11]string{"sun ", "moon", "merc", "ven ", "mars", "jup ", "sat ", "uran", "nept", "plut", "uzel"}
	namefilebd := "bw.txt" // файл базы данных удачных событий
	massiv := Makebd(namefilebd)
	step = 0
	itogo = 0 //счетчик найденного
	for massiv[step].Vrem != 0 {  
		m2 = massiv[step].Vrem
		asp2 := Makeasp(m2)
		//в оригинале sob = massiv[step].Sobyt
		sob = "Удачное событие"
		//сюда пойдут формулы
		pechatasp(20, asp2)  //печать 20 аспектов в порядке возрастания орбиса (для 2 теста)
		merc_nept(asp2)			
		step = step + 1
	}
	//fmt.Println()
	fmt.Println("\nСочетание: самый точный угол в аспекте с Нептуном, \nа самый точный межпланетный аспект - это Меркурий с планетой из таблицы ниже ")
	i := 0
	for i < 11{
		fmt.Println(names[i], summ[i])
		i = i + 1
	}	
	fmt.Println("найдено ", itogo, "из", step)
}
var m2 int
var step int
var sob string
func pechatzagolovka () {
	t := int64(m2) - 6*3600 //6*3600 - поправка на разницу с Гринв
	ut := time.Unix(t, 0)	
	fmt.Println("\n",step, ut.Format("02 Jan 2006 15:04:05"), sob)   
}
func pechatasp(kolichestvo int, dim[]Aspects ){ 
	i := 0
	for i < kolichestvo && i < len(dim){  
		g := int(dim[i].Orb)
		mi := int(math.Round(60 * (dim[i].Orb - float64(g))))	
		orbis := strconv.Itoa(g) + "°" + strconv.Itoa(mi) + "'"	
		hoz1 := "(" + dim[i].Upr1 + ")"
		hoz2 := "(" + dim[i].Upr2 + ")"
		fmt.Println(i, dim[i].Plnt1, hoz1, dim[i].Vid, dim[i].Plnt2, hoz2, orbis, dim[i].Shod)
		i = i + 1
	}		
}
// - - - - - - - - - - - - - - - - - - - - формулы - - - - - - - - - - - - - - 
func merc_nept(mass []Aspects) {	
	// ищем аспект Меркурия с x   	
	pp := "abv"
	if mass[0].Plnt1 == "merc"{		
		pp = mass[0].Plnt2
	}	
	if mass[0].Plnt2 == "merc"{		
		pp = mass[0].Plnt1
	}
	if pp == "abv" {
		return 
	} //Берем первый попавшийся аспект угла и смотрим, какую планету аспектирует угол
	i := 0	
	for i < len(mass){
		p1 := "xxx"
		p2 := "yyy"
		if mass[i].Plnt1 == "Asc" ||mass[i].Plnt1 == "MC"{
			p2 = mass[i].Plnt2
			if p2[:2] != "ne"{
				return 
			}
		}		
		if p2[:2] == "ne" ||p1[:2] == "ne"{
			pechatzagolovka()
			itogo = itogo + 1
			g := int(mass[i].Orb)
			mi := int(math.Round(60 * (mass[i].Orb - float64(g))))	
			orbis := strconv.Itoa(g) + "°" + strconv.Itoa(mi) + "'"	
			fmt.Println(i, mass[i].Plnt1, mass[i].Vid, mass[i].Plnt2, orbis, mass[i].Shod)
			pechatasp(1, mass)
			j := 0
			for j < 11 {
				if names[j][:2] == pp[:2] {
					summ[j] = summ[j] + 1
				}
				j = j + 1
			}
			return 
	        }
		i = i + 1		
	}		
	return 
}
//- - - - - - - - - конец искомой формулы - - - - - - - - - - - - - - - - - - - - -
// - - - - - - - - все, что ниже - можно вынести в отдельный файл - - - - - - - - - 
const dlinaEfm int = 43825 //длина массива структур 33  
const dlinabd int = 600
const minoraspects int = 0 // 0 - только мажорные аспекты
const separaspects int = 0 // 0 - не использовать сепарационные аспекты
const dlinaplanets int = 13// равна 13 когда отключен расчет жребиев
//  dlinaplanets int = 15 - когда  raschetjreb(planets) включена
const dlinaasp int = 300

func Makeasp(m2 int) []Aspects{
	ind, kk := poiskVrem(m2) //ind - индекс в efm, kk -коэф интерполяции
	zvtime := interpolstartime(ind, kk)		
	kus := raschetkuspidov(zvtime)		
	kus = alignkusp(kus)		
	planets := fillplanets(kus, ind, kk)		
	planets = plusUgly(planets, kus, zvtime) 
	//planets = raschetjreb(planets)
	asp := raschetasp(planets)		
	fillradix()
	asp2 := raschetasp2(planets, asp)
	return asp2
}
// создание массива с эфемеридами
func Makeefm() { //    читает эфемериды из файла и записывает в массив структур
	file, err := os.Open("file.txt")//количество строк в файле должно быть равно dlinaEfm)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	i := 0 
	for scanner.Scan() {
		a := scanner.Text()
		efm[i] = strEfemer(a)
		i = i + 1
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
func strEfemer(a string) Stroka {
	b := strings.Split(a, ";")
	daytime := b[0] + " " + b[1]	
	t, err := time.Parse("2.1.2006 15:04:05", daytime)
	if err != nil {
		panic(err)
	}
	var c Stroka
	c.vrem = int(t.Unix()) //время в секундах с начала эпохи Юникс (01.01.1970)
	c.startime = strToFloat(b[2])    
	c.sun = strToFloat(b[3])
	c.moon = strToFloat(b[4])
	c.merc = strToFloat(b[5])
	c.ven = strToFloat(b[6])
	c.mars = strToFloat(b[7])
	c.jup = strToFloat(b[8])
	c.sat = strToFloat(b[9])
	c.uran = strToFloat(b[10])
	c.nep = strToFloat(b[11])
	c.plu = strToFloat(b[12])
	c.uzel = strToFloat(b[13])
	return c
}

var efm [dlinaEfm]Stroka
func strToFloat(x string) float64 {
	var s float64
	s, err := strconv.ParseFloat(x, 64)
	if err != nil {
		panic(err)
	}
	return s
}

type Stroka struct {
	vrem     int //тип int для совместимости с функцией sort.SearchInts
	startime float64
	sun      float64
	moon     float64
	merc     float64
	ven      float64
	mars     float64
	jup      float64
	sat      float64
	uran     float64
	nep      float64
	plu      float64
	uzel     float64	
}
func poiskVrem(m2 int) (i int, kk float64) {  //поиск индекса (строки) в массиве efm, с наиболее близким временем к заданному
	h := m2 / 3600 
	ms := int(h * 3600)   //время в секундах, которое надо найти в массиве
	hf := float64(h)
	m2f := float64(m2)
	kk = m2f/3600.0 - hf //коэффициент корреляции
	//ищем время ms в массиве efm
	i = sort.Search(len(efm), func(i int) bool { return efm[i].vrem >= ms })  
	return i, kk 
}
var zvtime float64
func interpolstartime(ind int, kk float64) float64{
	st1 := efm[ind].startime 
	st2 := efm[ind+1].startime  
	ds := st2 - st1
	if ds < 0 {
		ds = ds + 24
	}
	dx := ds * kk
	zvtime = st1 + dx
	if zvtime > 24 {
		zvtime = zvtime - 24
	}
	return zvtime
}
func radians(grad float64) float64 {
	var rad float64 = grad * math.Pi / 180
	return rad
}
func degrees(rad float64) float64 {
	var grad float64 = rad * 180 / math.Pi
	return grad
}
func raschetkuspidov(zvtime float64) [12]float64 {
	stm := 360 * zvtime / 24
	ramc := radians(stm)
	e := radians(23.44166)
	mc := math.Atan(math.Tan(ramc) / math.Cos(e))
	k10 := degrees(mc)
	if math.Abs(stm-k10) > 10 {
		k10 = k10 + 180
	}
	if math.Abs(stm-k10) > 10 {  //повторно
		k10 = k10 + 180
	}
	mc = radians(k10)
	// расчет куспидов по системе Коха
	lat := radians(43.25)                        //географ. широта Алматы
	dec := math.Asin(math.Sin(mc) * math.Sin(e)) //склонение МС
	oamc1 := math.Tan(dec) * math.Tan(lat)       //наклонное восхождение МС
	oamcr := ramc - math.Asin(oamc1)             //в радианах
	oamc := degrees(oamcr)                       //в градусах
	dx := (stm + 90 - oamc) / 3                  //интервал между куспидами
	if dx < 0 {
		dx = dx + 360
	}
	h11 := oamc + dx - 90 //в градусах
	h12 := h11 + dx
	h1 := h12 + dx
	h2 := h1 + dx
	h3 := h2 + dx
	h11r := radians(h11) //в радианах
	h12r := radians(h12)
	h1r := radians(h1)
	h2r := radians(h2)
	h3r := radians(h3)
	//- - - - - - - - - - далее положение куспидов на эклиптике
	ts := math.Tan(lat) * math.Sin(e)
	e11r := math.Atan((-ts - math.Sin(h11r)*math.Cos(e)) / math.Cos(h11r))
	e11 := 90 - degrees(e11r)
	e12r := math.Atan((-ts - math.Sin(h12r)*math.Cos(e)) / math.Cos(h12r))
	e12 := 90 - degrees(e12r)
	e1r := math.Atan((-ts - math.Sin(h1r)*math.Cos(e)) / math.Cos(h1r))
	e1 := 90 - degrees(e1r)
	e2r := math.Atan((-ts - math.Sin(h2r)*math.Cos(e)) / math.Cos(h2r))
	e2 := 90 - degrees(e2r)
	e3r := math.Atan((-ts - math.Sin(h3r)*math.Cos(e)) / math.Cos(h3r))
	e3 := 90 - degrees(e3r)
	kus := [12]float64{e1, e2, e3, 0, 0, 0, 0, 0, 0, k10, e11, e12}
	return kus
}

func checkkusp(ku1 float64, ku2 float64) float64 { //функция для выравнивания куспидов
	razn := ku2 - ku1
	if razn >= -180 && ku2 <= ku1 {
		ku2 = ku2 + 180
	}
	return ku2
}
func checkkusp2(ku1 float64, ku2 float64) float64 { //функция для выравнивания куспидов
	razn := ku2 - ku1                           //второй вариант
	if razn >= -180 {
		if ku2 < ku1{
			ku2 = ku2 + 180
		}
	}
	return ku2
}
func alignkusp(kus [12]float64) [12]float64 {
	n := 9
	n6 := 0
	n2 := 9
	for n != 4 {		
		n2 = n + 1
		if n2 == 12 {
			n2 = 0
		}
		kus[n2] = checkkusp2(kus[n], kus[n2])
		n6 = n + 6
		if n6 > 11 {
			n6 = n6 - 12
		}
		kus[n6] = kus[n] + 180
		if kus[n6] > 360 {
			kus[n6] = kus[n6] - 360
		}
		n = n + 1
		if n == 12 {
			n = 0
		}		
	}
	return kus
}
func plusUgly(planets [15]plntdata, kus [12]float64, zvtime float64) [15]plntdata{
	zvtime2 := zvtime + 0.1 // 0.1 часа (6 минут) Звездного Времени
	asc1 := kus[0]
	mc1 := kus[9]
	kus2 := raschetkuspidov(zvtime2)    //
	asc2 := checkkusp(kus2[9], kus2[0]) //выравнивание Асц2
	mc2 := kus2[9]
	//считаем скорость в угловых секундах в час
	xmc := mc2 - mc1 //пробег МС в градусах
	if xmc < 0 {
		xmc = xmc + 360
	}
	vmc := 3600 * xmc / (36.0 / 361.0) // 36 = 360 * 0.1
	//360/361 - это соотношение между локальным временем и звездным
	xasc := asc2 - asc1
	if xasc < 0 {
		xasc = xasc + 360
	}
	var vasc float64 = 3600 * xasc / (36.0 / 361.0) // 36 = 360 * 0.1
	planets[0].name = "Asc"
	planets [0].koord = kus[0]
	planets [0].skor = vasc
	planets[1].name = "MC"
	planets [1].koord = kus[9]
	planets [1].skor = vmc
	return planets
}
func getkoord(ind int, kk float64, name string) (float64, float64){
	k1 := reflect.ValueOf(efm[ind]).FieldByName(name).Float()  
	k2 := reflect.ValueOf(efm[ind+1]).FieldByName(name).Float() 
	ds := k2 - k1  //  k2 - k1
	if ds < -300 {  //в случае пересечения О град Овна
		ds = ds + 360
		}
	v := 3600.0 * ds //скорость планеты, угловых секунд в час 
	var kx float64
	if ds < 0 {
		kx = k2 - ds * kk
		} else {
		kx = k1 + ds * kk
		}
	if kx >= 360 {
		kx = kx - 360.0
		}
	return kx, v  //координата и скорость планеты
}
type plntdata struct {
	name string
	koord float64
	skor float64	
	upr string  //для Узлов в этом поле пока ничего не пишется
}

var nazv[11]string
func fillplanets(kus [12]float64, ind int, kk float64) [15]plntdata{
	var planets [15]plntdata
	var kus2 [12]int
	i := 0
	for i < 12 {
		kus2[i] = int(kus[i]/30) + 1 //- - - - - - - в каких знаках куспиды домов
		i = i + 1
	}
	nazv = [11]string{"sun", "moon", "merc", "ven", "mars", "jup", "sat", "uran", "nep", "plu", "uzel"}
	z1 := [11]int{5, 4, 3, 2, 1, 9, 10, 11, 12, 8, 0} //- - - - - знаки под управл соотв планет
	z2 := [11]int{0, 0, 6, 7, 8, 12, 11, 0, 0, 0, 0}  //- - - - - -вторые знаки под управл соотв планет
	i = 2 // первые 2 строчки для Асц и МС, заполняются в plusUgly
	for i < 13 {
		planets[i].name = nazv[i-2]
		kx, v := getkoord(ind, kk, nazv[i-2])
		planets [i].koord = kx
		planets[i].skor = v
		n := 0
		for n < 12 {
			if kus2[n] == z1[i-2] || kus2[n] == z2[i-2] {
				x := strconv.Itoa(n + 1)
				if len(planets[i].upr) > 0 {
					x = ", " + x
				}
				planets[i].upr = planets[i].upr + x
			}
			n = n + 1
		}
		i = i + 1
	}
	return planets
}

func klassifasp(ugol float64) float64 {
	major := [5]float64{0.0, 60.0, 90.0, 120.0, 180.0}
	orbmain := 11.0
	minor := [4]float64{30, 45, 135, 150}
	qu := -1.0
	i := 0
	for i < len(major) {
		d := math.Abs(ugol - major[i])
		if d < orbmain {
			qu = major[i]
			return qu
		}
		i = i + 1
	}
	if minoraspects == 0 {
		return qu
	}
	orbminor := 2.0
	i = 0
	for i < len(minor) {
		d := math.Abs(ugol - minor[i])
		if d < orbminor {
			qu = minor[i]
			return qu
		}
		i = i + 1
	}
	return qu
}
// p1, p2 - координаты планет; v1, v2 - их скорости; u - угол между ними; 
// vid - вид аспекта (0, 60, 90 и т.п. в градусах)
func checkshod(p1, p2, v1, v2, u, vid float64) string {
	if vid == -1 {
		return "sep"
	}
	pr := radians(p1 - p2)
	pp := math.Sin(pr)
	var applic string
	if v1 > v2 {
		if pp < 0 {
			if u < vid {
				applic = "sep"
			} else {
				applic = "app"
			}
		} else {
			if u < vid {
				applic = "app"
			} else {
				applic = "sep"
			}
		}
	} else {
		if pp < 0 {
			if u < vid {
				applic = "app"
			} else {
				applic = "sep"
			}
		} else {
			if u < vid {
				applic = "sep"
			} else {
				applic = "app"
			}
		}
	}
	return applic
}		

type Aspects struct {
	Plnt1 string
	Upr1  string
	Vid   float64 	
	Plnt2 string
	Upr2  string
	Shod  string
	Orb   float64
}

func raschetasp(planets [15]plntdata) [dlinaasp]Aspects{
	i := 0
	z := 0
	var dim [dlinaasp]Aspects
	//dlin := len(planets) - 2  -2 когда отключен расчет жребиев
	for i < dlinaplanets {
		p1 := planets[i].koord
		hoz1 := planets[i].upr
		im1 := planets[i].name
		v1 := planets[i].skor
		j := i + 1
		for j < dlinaplanets {
			if i == j {
				j = j + 1
				continue
			}
			p2 := planets[j].koord
			u := math.Abs(p1 - p2)
			if u > 180 {
				u = 360.0 - u
			}
			vidasp := klassifasp(u)
			if vidasp == -1 {
				j = j + 1
				continue
			}
			
			v2 := planets[j].skor
			shod := checkshod(p1, p2, v1, v2, u, vidasp)
			if shod == "sep" && separaspects == 0{
				j = j + 1
				continue
			}
			hoz2 := planets[j].upr
			im2 := planets[j].name
			//- - - - - - - - - - - - запись данных в массив структур asp
			dim[z].Plnt1 = im1
			dim[z].Upr1 = hoz1
			//dim[z].ugol = u
			dim[z].Vid = vidasp			
			dim[z].Orb = math.Abs(u - vidasp)
			dim[z].Plnt2 = im2
			dim[z].Upr2 = hoz2
			dim[z].Shod = shod						
			z = z + 1
			j = j + 1			
		}
		i = i + 1
	}
	return dim
}

var radix [11]plntdata
func fillradix() {
	rnazv := [11]string{"sun.n", "moon.n", "merc.n", "ven.n", "mars.n", "jup.n", "sat.n", "uran.n", "nep.n", "plu.n", "uzel.n"}
	rKoords := [11]float64{96.8, 212.7, 99.8, 63.445, 132.192, 347.743, 98.12, 203.663, 247.5, 184.121, 259.358}
	rhoz := [11]string{"4", "2, 3", "1, 5", "12", "6", "7, 11", "9, 10", "10", "11", "6", "0"}
	i := 0
	for i < 11 {
		radix[i].name = rnazv[i]
		radix[i].koord = rKoords[i]
		radix[i].skor = 0
		radix[i].upr = rhoz[i]
		i = i + 1
	}	
}
func raschetasp2(planets [15]plntdata, asp [dlinaasp]Aspects) []Aspects{
	z := 0
	d := len(asp) - 1
	i := 0
	for asp[z].Plnt1 != "" && z < d {
		z = z + 1
	}	
	//dlin := len(planets) - 2    -2 когда отключен расчет жребиев
	for i < dlinaplanets {
		p1 := planets[i].koord
		hoz1 := planets[i].upr
		im1 := planets[i].name
		v1 := planets[i].skor
		j := 0  // здесь надо начинать с нуля
		for j < 11 { //длина radix
			p2 := radix[j].koord
			u := math.Abs(p1 - p2)
			if u > 180 {
				u = 360.0 - u
			}
			vidasp := klassifasp(u)
			if vidasp == -1 {
				j = j + 1
				continue
			}
			v2 := 0.0
			shod := checkshod(p1, p2, v1, v2, u, vidasp)
			if shod == "sep" && separaspects == 0{
				j = j + 1
				continue
			}
			hoz2 := radix[j].upr
			im2 := radix[j].name
			//- - - - - - - - - - запись данных в массив структур asp
			asp[z].Plnt1 = im1
			asp[z].Upr1 = hoz1
			//asp[z].ugol = u
			asp[z].Vid = vidasp			
			asp[z].Orb = math.Abs(u - vidasp)
			asp[z].Plnt2 = im2
			asp[z].Upr2 = hoz2
			asp[z].Shod = shod			
			z = z + 1				
			j = j + 1			
		}
		i = i + 1
	}
	srez := asp[:z] // - - - - - берем только непустые строки - - - - - 
	sort.Slice(srez, func(i, j int) bool {
		return srez[i].Orb < srez[j].Orb
	})	
	return srez
}
func raschetjreb(planets [15]plntdata) [15]plntdata{
	asc := planets[11].koord
	sun := planets[0].koord
	moon := planets[1].koord
	var den float64 = 0
	das := asc - sun //определение дня или ночи
	if das < -180 {
		das = 360 + das
	}
	if das > 0 {
		den = 1 //это день
	} else {
		den = -1 //это ночь
	}
	fortune := asc + den*(moon-sun)
	tduha := asc + den*(sun-moon)
	if fortune > 360 {
		fortune = fortune - 360
	}
	if fortune < 0 {
		fortune = fortune + 360
	}
	if tduha > 360 {
		tduha = tduha - 360
	}
	if tduha < 0 {
		tduha = tduha + 360
	}
	//fmt.Println("Hello", fortune, tduha)
	planets[13].name = "Fortune"
	planets[13].koord = fortune
	planets[13].skor = planets[11].skor // скорость этих точек практически равна скорости Асц
	planets[14].name = "tDuha"
	planets[14].koord = tduha
	planets[14].skor = planets[11].skor
	return planets
}

// массивы bw, bp, ba, bf, bpf должны быть одинаковой длины, чтобы их можно было использовать как аргумент для одной и той же функции
type Strbazdan struct {
	Vrem     int //тип int для совместимости с функцией sort.SearchInts
	Sobyt string
}

func strbd(a string) Strbazdan {  //для функции makeBd
	b := strings.Split(a, ";")
	if len(b[1]) < 6 {  //для случаев, когда наподобие "21.02.2020 6:50"
		b[1] = b[1] + ":00"
	}
	daytime := b[0] + " " + b[1]
	t, err := time.Parse("2.1.2006 15:04:05", daytime)	
	if err != nil {
		panic(err)
	}
	var c Strbazdan
	c.Vrem = int(t.Unix())  //время в секундах с начала эпохи Юникс (01.01.1970)
	c.Sobyt = b[2]		       	
	return c
}
func Makebd(filebd string) [dlinabd]Strbazdan{ //    читает эфемериды из файла и записывает в массив структур
	file, err := os.Open(filebd)//количество строк в файле должно быть не более dlinabd
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var massiv [dlinabd]Strbazdan
	scanner := bufio.NewScanner(file)
	i := 0 
	for scanner.Scan() {
		a := scanner.Text()
		massiv[i] = strbd(a)
		i = i + 1
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return massiv
}	


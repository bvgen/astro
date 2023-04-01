# astro
Python vs Golang
Поиск определенных астрологических сочетаний в списке событий или в диапазоне дат
There are 2 versions of this program - in Python and in Golang, single-threaded.
Since there are implementations in different languages, we can compare their performance.
Depending on launch option, Golang can be faster than Python up to 400 times faster (40 000%).
This code searches for given combinations of planets in a list of events or a range of dates.
First, it loads the ephemeris, then determines the coordinates of the planets for a certain time, then calculates the aspects (angles) between the planets.
In fact, a horoscope is built for each event.
Secondly, in the resulting list of aspects, a user-specified combination of planets is searched.
Ephemeris obtained from Zet program, https://astrozet.net 
Which has the function of creating ephemeris tables.
Since there are implementations in different languages, we can compare their performance.

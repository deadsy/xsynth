#!/usr/bin/python

import math


THRESHOLD = 0.03

def main():

  vs = 1.0
  vf = 0.0

  n = 15
  r = math.exp(math.log(THRESHOLD) / n)
  k = 1.0 - r

  val = vs

  for i in range(n + 1):
    print('%d: %f' % (i, val))
    val += k * (vf - val)

  for i in range(n + 1):
    ri = math.pow(r, i)
    val = vf + (vs - vf) * ri
    print('%d: %f' % (i, val))


main()




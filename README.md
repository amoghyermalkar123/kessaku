Kessaku is a worker pool package aimed at simplicity and solving issues in golang where high-performance systems run into bottleneck issues caused by the cost of creating goroutines with increasing stack-size.

Kessaku lets you create an arbitrary pool of workers that can be re-used, resulting in lesser and un-predictable `morestack` calls by ever-increasing goroutines.
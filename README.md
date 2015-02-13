# Geo Index

Geo Index library

## Overview

Splits the earth surface in a grid. At each cell we can store data, such as list of points, count of points, etc. It can do KNearest and Within queries.

### Demo

http://go-geoindex.appspot.com/static/nearest.html - Click to select the nearest points.

http://go-geoindex.appspot.com/static/cluster.html - A map with 100K points around the world. Zoom in and out to cluster. 

### API

```go
    type Driver struct {
        lat float64
        lon float64
        id string
        canAcceptJobs bool
    }

    // Implement Point interface
    func (d *Driver) Lat() { return d.lat }
    func (d *Driver) Lon() { return d.lat }
    func (d *Driver) Id() { return d.id }

    // create points index with resolution (cell size) 0.5 km
    index := NewPointsIndex(Km(0.5))

    // Adds a point in the index, if a point with the same id exists it's removed and the new one is added
    index.Add(&Driver{id1, lat, lng, true})
    index.Add(&Driver{id2, lat, lng, false})

    // Removes a point from the index by id
    index.Remove(&Driver{id1, lat, lng, true})

    // get the k-nearest points to a point, within some distance
    points := index.KNearest(&GeoPoint{id, lat, lng}, 5, Km(5), func(p Point) bool {
        return p.(* Driver).canAcceptJobs
    })

    // get the points within a range on the map
    points := index.Within(topLeftPoint, bottomRightPoint)
```

### Index types

There are several index types

```go
    NewPointsIndex(Km(0.5)) // Creates index that maintains points
    NewExpiringPointsIndex(Km(0.5), Minutes(5)) // Creates index that expires the points after some interval
    NewCountIndex(Km(0.5)) // Creates index that maintains counts of the points in each cell
    NewExpiringCountIndex(Km(0.5), Minutes(15)) // Creates index that maintains expiring count
    NewClusteringIndex() // index that clusters the points at different zoom levels, so we can create maps
    NewExpiringClusteringIndex(Minutes(15)) // index that clusters and expires the points at different zoom levels
                                            // so we can create real time maps of customer request, etc in the driver app
```

### Performance Benchmarks

    BenchmarkClusterIndexAdd	                500000	      5965 ns/op
    BenchmarkClusterIndexWithinStreet	        200000	     10205 ns/op
    BenchmarkClusterIndexWithinCity	            100000	     19408 ns/op
    BenchmarkClusterIndexWithinWorld	        50000	     32226 ns/op

    BenchmarkExpiringClusterIndexAdd	        500000	      5250 ns/op
    BenchmarkExpiringClusterIndexWithinStreet	200000	     14887 ns/op
    BenchmarkExpiringClusterIndexWithinCity	    100000	     21920 ns/op
    BenchmarkExpiringClusterIndexWithinWorld	50000	     32737 ns/op

    BenchmarkCountIndexAdd	                    1000000	      1327 ns/op
    BenchmarkCountIndexWithin	                200000	     12419 ns/op

    BenchmarkExpiringCountIndexAdd	            1000000	      2273 ns/op
    BenchmarkExpiringCountIndexWithin	        100000	     16535 ns/op

    BenchmarkPointIndexWithin	                200000	      9288 ns/op
    BenchmarkPointIndexAdd	                    1000000	      2174 ns/op
    BenchmarkPointIndexKNearest	                100000	     15137 ns/op

    BenchmarkExpiringPointIndexAdd	            1000000	      2746 ns/op
    BenchmarkExpiringPointIndexKNearest	        100000	     17689 ns/op
    BenchmarkExpiringPointIndexWithin	        200000	      9741 ns/op


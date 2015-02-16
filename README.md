# Geo Index

Geo Index library

## Overview

Splits the earth surface in a grid. At each cell we can store data, such as list of points, count of points, etc. It can do KNearest and Range queries.

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
    points := index.Range(topLeftPoint, bottomRightPoint)
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

    BenchmarkClusterIndexAdd                    300000         5068 ns/op
    BenchmarkClusterIndexStreetRange            100000        23611 ns/op
    BenchmarkClusterIndexCityRange              30000         47462 ns/op
    BenchmarkClusterIndexEuropeRange            50000         32509 ns/op

    BenchmarkExpiringClusterIndexAdd            200000         6431 ns/op
    BenchmarkExpiringClusterIndexStreetRange    50000         27730 ns/op
    BenchmarkExpiringClusterIndexCityRange      20000         66127 ns/op
    BenchmarkExpiringClusterIndexEuropeRange    30000         39111 ns/op

    BenchmarkCountIndexAdd                      1000000        2210 ns/op
    BenchmarkCountIndexRange                    30000         63263 ns/op    

    BenchmarkExpiringCountIndexAdd              300000         4191 ns/op
    BenchmarkExpiringCountIndexRange            30000         59754 ns/op
    
    BenchmarkPointIndexAdd                      500000         3981 ns/op
    BenchmarkPointIndexRange                    50000         30816 ns/op
    BenchmarkPointIndexKNearest                 50000         22854 ns/op

    BenchmarkExpiringPointIndexAdd              200000          5129 ns/op
    BenchmarkExpiringPointIndexKNearest         100000         16598 ns/op
    BenchmarkExpiringPointIndexRange            100000         18911 ns/op

window.BENCHMARK_DATA = {
  "lastUpdate": 1732579326551,
  "repoUrl": "https://github.com/BeryJu/gravity",
  "entries": {
    "Gravity Benchmark": [
      {
        "commit": {
          "author": {
            "name": "BeryJu",
            "username": "BeryJu"
          },
          "committer": {
            "name": "BeryJu",
            "username": "BeryJu"
          },
          "id": "f97b32e08812bdd7b71d4a19594718aff8e458f5",
          "message": "dns: start adding benchmarks for performance regressions",
          "timestamp": "2024-11-25T17:19:56Z",
          "url": "https://github.com/BeryJu/gravity/pull/1325/commits/f97b32e08812bdd7b71d4a19594718aff8e458f5"
        },
        "date": 1732579326248,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkRoleDNS_Memory",
            "value": 9529,
            "unit": "ns/op\t    4944 B/op\t      66 allocs/op",
            "extra": "126408 times\n4 procs"
          },
          {
            "name": "BenchmarkRoleDNS_Memory - ns/op",
            "value": 9529,
            "unit": "ns/op",
            "extra": "126408 times\n4 procs"
          },
          {
            "name": "BenchmarkRoleDNS_Memory - B/op",
            "value": 4944,
            "unit": "B/op",
            "extra": "126408 times\n4 procs"
          },
          {
            "name": "BenchmarkRoleDNS_Memory - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "126408 times\n4 procs"
          }
        ]
      }
    ]
  }
}
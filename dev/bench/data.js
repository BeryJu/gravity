window.BENCHMARK_DATA = {
  "lastUpdate": 1732580038260,
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
      },
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
          "id": "e582d43cbb87185e17bd70a5012c99ac45e45c31",
          "message": "dns: start adding benchmarks for performance regressions",
          "timestamp": "2024-11-25T17:19:56Z",
          "url": "https://github.com/BeryJu/gravity/pull/1325/commits/e582d43cbb87185e17bd70a5012c99ac45e45c31"
        },
        "date": 1732580037391,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkRoleDNS_Etcd",
            "value": 720707,
            "unit": "ns/op\t   29921 B/op\t     413 allocs/op",
            "extra": "1572 times\n4 procs"
          },
          {
            "name": "BenchmarkRoleDNS_Etcd - ns/op",
            "value": 720707,
            "unit": "ns/op",
            "extra": "1572 times\n4 procs"
          },
          {
            "name": "BenchmarkRoleDNS_Etcd - B/op",
            "value": 29921,
            "unit": "B/op",
            "extra": "1572 times\n4 procs"
          },
          {
            "name": "BenchmarkRoleDNS_Etcd - allocs/op",
            "value": 413,
            "unit": "allocs/op",
            "extra": "1572 times\n4 procs"
          },
          {
            "name": "BenchmarkRoleDNS_Memory",
            "value": 9660,
            "unit": "ns/op\t    4944 B/op\t      66 allocs/op",
            "extra": "116751 times\n4 procs"
          },
          {
            "name": "BenchmarkRoleDNS_Memory - ns/op",
            "value": 9660,
            "unit": "ns/op",
            "extra": "116751 times\n4 procs"
          },
          {
            "name": "BenchmarkRoleDNS_Memory - B/op",
            "value": 4944,
            "unit": "B/op",
            "extra": "116751 times\n4 procs"
          },
          {
            "name": "BenchmarkRoleDNS_Memory - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "116751 times\n4 procs"
          }
        ]
      }
    ]
  }
}
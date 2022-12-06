# BPlusTree

B+ tree in Go.

## 1. Insertion

```txt
=== RUN   TestInsertNode
    bplustree_test.go:133: tree after insert key: 3
          (3)

    bplustree_test.go:133: tree after insert key: 1
          (1,3)

    bplustree_test.go:133: tree after insert key: 2
              (3)
          (1,2)    (3)

    bplustree_test.go:133: tree after insert key: 4
               (3)
          (1,2)    (3,4)

    bplustree_test.go:133: tree after insert key: 5
                  (3,5)
          (1,2)    (3,4)    (5)

    bplustree_test.go:133: tree after insert key: 6
                   (3,5)
          (1,2)    (3,4)    (5,6)

    bplustree_test.go:133: tree after insert key: 7
                       (5)
               (3)             (7)
          (1,2)    (3,4)    (5,6)    (7)

    bplustree_test.go:133: tree after insert key: 20
                         (5)
               (3)               (7)
          (1,2)    (3,4)    (5,6)    (7,20)

    bplustree_test.go:133: tree after insert key: 18
                             (5)
               (3)                 (7,20)
          (1,2)    (3,4)    (5,6)    (7,18)    (20)

    bplustree_test.go:133: tree after insert key: 19
                               (5,19)
               (3)               (7)              (20)
          (1,2)    (3,4)    (5,6)    (7,18)    (19)    (20)
```

## 2. Deletion

```txt
=== RUN   TestDeleteInternalNodeBorrowRight
    bplustree_test.go:124: tree after insert key: 1
          (1)
    bplustree_test.go:124: tree after insert key: 2
          (1,4)
    bplustree_test.go:124: tree after insert key: 3
          (1,4,7)
    bplustree_test.go:124: tree after insert key: 4
                (7)
          (1,4)    (7,10)
    bplustree_test.go:124: tree after insert key: 5
                 (7)
          (1,4)    (7,10,13)
    bplustree_test.go:124: tree after insert key: 6
                    (7,13)
          (1,4)    (7,10)    (13,16)
    bplustree_test.go:124: tree after insert key: 7
                     (7,13)
          (1,4)    (7,10)    (13,16,19)
    bplustree_test.go:124: tree after insert key: 8
                        (7,13,19)
          (1,4)    (7,10)    (13,16)    (19,22)
    bplustree_test.go:124: tree after insert key: 9
                         (7,13,19)
          (1,4)    (7,10)    (13,16)    (19,22,25)
    bplustree_test.go:124: tree after insert key: 10
                                (19)
                    (7,13)                     (25)
          (1,4)    (7,10)    (13,16)    (19,22)    (25,28)
    bplustree_test.go:264: b tree:
                                     (19)
                        (7,11,13)                         (25)
          (1,4)    (7,10)    (11,12)    (13,16)    (19,22)    (25,28)

    bplustree_test.go:269: b tree after deleting key 19:
                                (13)
                    (7,11)                     (19)
          (1,4)    (7,10)    (11,12)    (13,16)    (25,28)
```

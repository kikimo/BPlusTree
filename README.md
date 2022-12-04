# BPlusTree

Go implementation of BPlusTree.

## Insertion

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

    bplustree_test.go:133: tree after insert key: 13
                                   (5,19)
               (3)                 (7,18)                (20)
          (1,2)    (3,4)    (5,6)    (7,13)    (18)    (19)    (20)

    bplustree_test.go:133: tree after insert key: 10
                                        (13)
                         (5)                               (19)
               (3)               (7)              (18)            (20)
          (1,2)    (3,4)    (5,6)    (7,10)    (13)    (18)    (19)    (20)

    bplustree_test.go:133: tree after insert key: 12
                                            (13)
                             (5)                                   (19)
               (3)                 (7,12)                (18)            (20)
          (1,2)    (3,4)    (5,6)    (7,10)    (12)    (13)    (18)    (19)    (20)

    bplustree_test.go:133: tree after insert key: 11
                                                (13)
                               (5,11)                                     (19)
               (3)               (7)              (12)            (18)            (20)
          (1,2)    (3,4)    (5,6)    (7,10)    (11)    (12)    (13)    (18)    (19)    (20)

    bplustree_test.go:133: tree after insert key: 17
                                                  (13)
                               (5,11)                                      (19)
               (3)               (7)              (12)             (18)             (20)
          (1,2)    (3,4)    (5,6)    (7,10)    (11)    (12)    (13,17)    (18)    (19)    (20)

    bplustree_test.go:133: tree after insert key: 16
                                                      (13)
                               (5,11)                                          (19)
               (3)               (7)              (12)                (17,18)                (20)
          (1,2)    (3,4)    (5,6)    (7,10)    (11)    (12)    (13,16)    (17)    (18)    (19)    (20)

    bplustree_test.go:133: tree after insert key: 14
                                                          (13)
                               (5,11)                                             (17,19)
               (3)               (7)              (12)             (16)             (18)            (20)
          (1,2)    (3,4)    (5,6)    (7,10)    (11)    (12)    (13,14)    (16)    (17)    (18)    (19)    (20)

    bplustree_test.go:133: tree after insert key: 15
                                                              (13)
                               (5,11)                                                 (17,19)
               (3)               (7)              (12)                (15,16)                (18)            (20)
          (1,2)    (3,4)    (5,6)    (7,10)    (11)    (12)    (13,14)    (15)    (16)    (17)    (18)    (19)    (20)

    bplustree_test.go:133: tree after insert key: 9
                                                                 (13)
                                   (5,11)                                                     (17,19)
               (3)                 (7,10)                (12)                (15,16)                (18)            (20)
          (1,2)    (3,4)    (5,6)    (7,9)    (10)    (11)    (12)    (13,14)    (15)    (16)    (17)    (18)    (19)    (20)

    bplustree_test.go:133: tree after insert key: 8
                                                                    (9,13)
                        (5)                             (11)                                       (17,19)
               (3)              (7)            (10)           (12)                (15,16)                (18)            (20)
          (1,2)    (3,4)    (5,6)    (7,8)    (9)    (10)    (11)    (12)    (13,14)    (15)    (16)    (17)    (18)    (19)    (20)

--- PASS: TestInsertNode (0.00s)
```

func TestIsFull(t *testing.T) {
    group := NewRunGroup(5)

    Convey("Add < max workers", t, func() {
        for i := 0; i < 4; i++ {
            group.Add(1)
            So(group.IsFull(), ShouldBeFalse) // HL
        }
    })
}

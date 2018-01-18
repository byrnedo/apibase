package protobuf

import "time"

func (m *NatsContext_Trail) SetDeadlineTs(ts time.Time) {
	s := ts.Unix()
	n := int32(ts.Nanosecond())

	m.Deadline = &s
	m.DeadlineNanos = &n
}

func (m *NatsContext_Trail) GetDeadlineTs() *time.Time {
	if m.Deadline == nil {
		return nil
	}

	dS := m.GetDeadline()
	dN := m.GetDeadlineNanos()

	t := time.Unix(dS, int64(dN))
	return &t
}

func (m *NatsContext_Trail) SetTimeTs(ts time.Time) {
	s := ts.Unix()
	n := int32(ts.Nanosecond())

	m.Time = &s
	m.TimeNanos = &n
}

func (m *NatsContext_Trail) GetTimeTs() *time.Time {
	if m.Time == nil {
		return nil
	}

	dS := m.GetTime()
	dN := m.GetTimeNanos()

	t := time.Unix(dS, int64(dN))
	return &t
}


func (m *NatsContext) GetLastTrail() *NatsContext_Trail {
	if len(m.Trail) < 1 {
		return nil
	}

	return m.Trail[len(m.Trail) - 1]
}
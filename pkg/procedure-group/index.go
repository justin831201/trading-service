package procedure_group

type procedure struct {
	execute   func() error
	interrupt func(error)
}

type ProcedureGroup struct {
	procedures []procedure
}

func (group *ProcedureGroup) Add(execute func() error, interrupt func(err error)) {
	group.procedures = append(group.procedures, procedure{
		execute:   execute,
		interrupt: interrupt,
	})
}

func (group *ProcedureGroup) Run() error {
	if len(group.procedures) == 0 {
		return nil
	}

	errors := make(chan error, len(group.procedures))
	for _, p := range group.procedures {
		go func(p procedure) {
			errors <- p.execute()
		}(p)
	}

	err := <-errors
	for _, p := range group.procedures {
		p.interrupt(err)
	}
	for i := 1; i < len(group.procedures); i++ {
		<-errors
	}
	return err
}

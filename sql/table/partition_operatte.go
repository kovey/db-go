package table

import (
	"strconv"
	"strings"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql/operator"
)

type PartitionOperate struct {
	partitionNames []string
	method         string
	table          string
	validation     string
	number         int
	define         *PartitionDefinition
	defines        []*PartitionDefinition
}

func NewPartitionOperate(names ...string) *PartitionOperate {
	return &PartitionOperate{partitionNames: names}
}

func (p *PartitionOperate) Drop() ksql.PartitionOperateInterface {
	p.method = "DROP"
	return p
}

func (p *PartitionOperate) Discard() ksql.PartitionOperateInterface {
	p.method = "DISCARD PARTITION"
	return p
}

func (p *PartitionOperate) Import() ksql.PartitionOperateInterface {
	p.method = "IMPORT PARTITION"
	return p
}

func (p *PartitionOperate) Truncate() ksql.PartitionOperateInterface {
	p.method = "TRUNCATE PARTITION"
	return p
}

func (p *PartitionOperate) Reorganize() ksql.PartitionDefinitionInterface {
	p.method = "REORGANIZE PARTITION"
	define := &PartitionDefinition{name: p.partitionNames[0], onlyBody: true}
	p.defines = append(p.defines, define)
	p.partitionNames = nil
	return define
}

func (p *PartitionOperate) Exchange(table, validation string) ksql.PartitionOperateInterface {
	p.method = "EXCHANGE PARTITION"
	p.table = table
	p.validation = validation
	return p
}

func (p *PartitionOperate) Analyze() ksql.PartitionOperateInterface {
	p.method = "ANALYZE PARTITION"
	return p
}

func (p *PartitionOperate) Check() ksql.PartitionOperateInterface {
	p.method = "CHECK PARTITION"
	return p
}

func (p *PartitionOperate) Optimize() ksql.PartitionOperateInterface {
	p.method = "OPTIMIZE PARTITION"
	return p
}

func (p *PartitionOperate) Rebuild() ksql.PartitionOperateInterface {
	p.method = "REBUILD PARTITION"
	return p
}

func (p *PartitionOperate) Repair() ksql.PartitionOperateInterface {
	p.method = "REPAIR PARTITION"
	return p
}

func (p *PartitionOperate) Remove() ksql.PartitionOperateInterface {
	p.method = "REMOVE PARTITIONING"
	return p
}

func (p *PartitionOperate) Coalesce(number int) ksql.PartitionOperateInterface {
	p.method = "COALESCE PARTITION"
	p.number = number
	return p
}

func (p *PartitionOperate) Add() ksql.PartitionDefinitionInterface {
	p.method = "ADD"
	p.define = &PartitionDefinition{name: p.partitionNames[0], onlyBody: true}
	p.partitionNames = nil
	return p.define
}

func (p *PartitionOperate) isAll() bool {
	for _, name := range p.partitionNames {
		if name == "ALL" {
			return true
		}
	}

	return false
}

func (p *PartitionOperate) Build(builder *strings.Builder) {
	builder.WriteString(p.method)
	if len(p.partitionNames) > 0 {
		builder.WriteString(" ")
		if p.isAll() {
			builder.WriteString("ALL")
		} else {
			for index, name := range p.partitionNames {
				if index > 0 {
					builder.WriteString(", ")
				}

				operator.Column(name, builder)
			}
		}
	}

	if p.number > 0 {
		builder.WriteString(" ")
		builder.WriteString(strconv.Itoa(p.number))
	}

	switch p.method {
	case "ADD":
		builder.WriteString("(")
		p.define.Build(builder)
		builder.WriteString(")")
	case "DISCARD PARTITION", "IMPORT PARTITION":
		builder.WriteString(" TABLESPACE")
	case "REORGANIZE PARTITION":
		builder.WriteString(" INTO (")
		for index, define := range p.defines {
			if index > 0 {
				builder.WriteString(", ")
			}

			define.Build(builder)
		}
		builder.WriteString(")")
	case "EXCHANGE PARTITION":
		builder.WriteString(" WITH TABLE ")
		operator.Column(p.table, builder)
		if p.validation != "" {
			builder.WriteString(" ")
			builder.WriteString(p.validation)
			builder.WriteString(" VALIDATION")
		}
	}
}

type PartitionOperates struct {
	operates []*PartitionOperate
}

func (p *PartitionOperates) Append(po *PartitionOperate) *PartitionOperates {
	p.operates = append(p.operates, po)
	return p
}

func (p *PartitionOperates) Build(builder *strings.Builder) {
	for _, po := range p.operates {
		builder.WriteString(" ")
		po.Build(builder)
	}
}

func (p *PartitionOperates) Empty() bool {
	return len(p.operates) == 0
}

package gospider

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"time"
)
var TABLE_FIELDS = [5]string{"time","label","url","action","respy"}
type ActionRecorder struct{
	db *sql.DB
	stmtPut *sql.Stmt
	stmtGet *sql.Stmt
	stmtDel *sql.Stmt
	label string
}

func (this *ActionRecorder) Init(cfg ActionRecordConfig) error{
	var err error
	this.db,err = sql.Open(cfg.Type,cfg.User+":"+cfg.PassWord+"@tcp("+cfg.Address+")/"+cfg.DB)
	if err!=nil{
		return err
	}
	err=this.db.Ping()
	if err != nil{
		return err
	}
	s:=fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(
		id BIGINT NOT NULL auto_increment,
		PRIMARY KEY (id)
	)ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci;`,cfg.Table)
	_,err=this.db.Exec(s)
	if err !=nil{
		return err
	}
	err = this.checkFeilds(this.db,cfg.Table)
	if err != nil{
		return err
	}
	fls:=""
	values:=""
	for _,v:=range TABLE_FIELDS{
		fls+=v+","
		values+="?,"
	}
	fls=strings.TrimRight(fls,",")
	values=strings.TrimRight(values,",")
	//time,label,url,content,failed,func,method,postdata
	this.stmtPut,err=this.db.Prepare(fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`,cfg.Table,fls,values))
	if err != nil{
		return err
	}
	this.stmtGet,err=this.db.Prepare(fmt.Sprintf(`SELECT id,action FROM %s WHERE label="%s" AND respy<=%d`,cfg.Table,this.label,cfg.MaxRespy))
	if err != nil{
		return err
	}
	this.stmtDel,err=this.db.Prepare(fmt.Sprintf(`DELETE FROM %s WHERE id=?`,cfg.Table))
	if err != nil{
		return err
	}
	gob.Register(Meta{})
	return nil
}

func (this *ActionRecorder) checkFeilds(db *sql.DB,tb string) error{
	fields:= make([]string,len(TABLE_FIELDS))
	copy(fields,TABLE_FIELDS[:])
	rows, err := db.Query("SELECT COLUMN_NAME FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='"+tb+"'")
	if err != nil{
		return err
	}
	defer rows.Close()
	var field string
	for rows.Next(){
		err:=rows.Scan(&field)
		if err != nil{
			return err
		}
		for k,v := range fields{
			if v == field{
				fields=append(fields[:k],fields[k+1:]...)
				break
			}
		}
	}
	for _,v:= range fields{
		switch v {
		case "label":
			_,err:=db.Exec(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s VARCHAR(64)",tb,v))
			if err != nil{
				return err
			}
		case "action":
			_,err:=db.Exec(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s BLOB",tb,v))
			if err != nil{
				return err
			}
		case "url":
			_,err:=db.Exec(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s VARCHAR(255)",tb,v))
			if err != nil{
				return err
			}
		case "respy":
			_,err:=db.Exec(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s INT",tb,v))
			if err != nil{
				return err
			}
		case "time":
			_,err:=db.Exec(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s DATETIME",tb,v))
			if err != nil{
				return err
			}
		}
	}
	return nil
}

func (this *ActionRecorder) Put(action Action) error{
	action.failCount=0
	var binary bytes.Buffer
	encoder := gob.NewEncoder(&binary)
	err:=encoder.Encode(action)
	b:=binary.Bytes()
	_,err=this.stmtPut.Exec(time.Now().Format("2006-01-02 15:04:05"),this.label,action.Url,b, action.Respy)
	if err != nil{
		return err
	}
	return nil
}

func (this *ActionRecorder) SetActionLabel(label string){
	this.label=label
}

func (this *ActionRecorder) GetActions() (actions []Action,err error){
	rows,err:=this.stmtGet.Query()
	if err != nil{
		return
	}
	defer rows.Close()
	var ids[]int
	for rows.Next(){
		var id int
		var content []byte

		err=rows.Scan(&id, &content)
		if err!=nil{
			return
		}
		var actiontemp Action
		decoder := gob.NewDecoder(bytes.NewBuffer(content))
		err=decoder.Decode(&actiontemp)
		if err!=nil{
			return
		}
		action:= actiontemp.UnsafeClone()
		action.Respy=actiontemp.Respy
		actions=append(actions,action)
		ids=append(ids,id)
	}
	for _,id:=range ids{
		_,err:=this.stmtDel.Exec(id)
		if err != nil{
			log.Fatal(err)
		}
	}
	return
}

func (this *ActionRecorder) Close(){
	this.stmtPut.Close()
	this.stmtGet.Close()
	this.stmtDel.Close()
	this.db.Close()
}
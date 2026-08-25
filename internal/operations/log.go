package operations

import("fmt";"sort";"strings";"time")
type LogEntry struct{ID,Level,Message,Actor string; At time.Time; Fields map[string]string}
type Log struct{entries []LogEntry}
func NewLog()*Log{return &Log{entries:[]LogEntry{}}}
func(l *Log) Add(e LogEntry)error{if e.ID==""||e.Message==""{return fmt.Errorf("log identity required")};if e.Level==""{e.Level="info"};if e.At.IsZero(){return fmt.Errorf("timestamp required")};e.Fields=clone(e.Fields);l.entries=append(l.entries,e);sort.SliceStable(l.entries,func(i,j int)bool{return l.entries[i].At.Before(l.entries[j].At)});return nil}
func(l *Log) List(level string)[]LogEntry{out:=[]LogEntry{};for _,e:=range l.entries{if level!=""&&e.Level!=level{continue};e.Fields=clone(e.Fields);out=append(out,e)};return out}
func(l *Log) Since(at time.Time)[]LogEntry{out:=[]LogEntry{};for _,e:=range l.entries{if e.At.Before(at){continue};out=append(out,e)};return out}
func(l *Log) ByActor(actor string)[]LogEntry{out:=[]LogEntry{};for _,e:=range l.entries{if e.Actor==actor{out=append(out,e)}};return out}
func(l *Log) Count()int{return len(l.entries)}
func(l *Log) Has(id string)bool{for _,e:=range l.entries{if e.ID==id{return true}};return false}
func(l *Log) Last() (LogEntry,bool){if len(l.entries)==0{return LogEntry{},false};return l.entries[len(l.entries)-1],true}
func(l *Log) Levels()[]string{m:=map[string]bool{};for _,e:=range l.entries{m[e.Level]=true};out:=[]string{};for x:=range m{out=append(out,x)};sort.Strings(out);return out}
func Format(e LogEntry)string{return fmt.Sprintf("%s [%s] %s: %s",e.At.Format(time.RFC3339),strings.ToUpper(e.Level),e.Actor,e.Message)}
func clone(in map[string]string)map[string]string{out:=map[string]string{};for k,v:=range in{out[k]=v};return out}

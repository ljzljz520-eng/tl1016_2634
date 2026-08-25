package intake

import("fmt";"strings")
type Field struct{Name,Value string; Required bool}
type Fields struct{values map[string]string}
func NewFields()*Fields{return &Fields{values:map[string]string{}}}
func(f *Fields) Set(name,value string)error{if strings.TrimSpace(name)==""{return fmt.Errorf("name required")};f.values[name]=strings.TrimSpace(value);return nil}
func(f *Fields) Get(name string)(string,bool){v,ok:=f.values[name];return v,ok}
func(f *Fields) Require(names ...string)[]string{out:=[]string{};for _,n:=range names{if strings.TrimSpace(f.values[n])==""{out=append(out,n)}};return out}
func(f *Fields) Delete(name string){delete(f.values,name)}
func(f *Fields) Count()int{return len(f.values)}
func(f *Fields) Empty()bool{return len(f.values)==0}
func(f *Fields) Keys()[]string{out:=[]string{};for k:=range f.values{out=append(out,k)};return out}
func(f *Fields) Clone()*Fields{n:=NewFields();for k,v:=range f.values{n.values[k]=v};return n}
func Join(values []string)string{return strings.Join(values,", ")}

名称	    特点与定位	    独立性	    方法和类的属性
util	通用的、与业务无关的，可以独立出来，可供其他项目使用	不调用任何业务相关的类  	方法通常是public static的，一般无类的属性，如果有，也是public static的
tool	与某些业务有关，通用性只限于某几个业务类之间	要调用某些业务相关的类	方法通常是public static的，一般无类的属性，如果有，也是public static的
service	与某一个业务有关，不是通用的	要调用某些业务相关的类	方法通常是public的，通常是通过接口去调用，一般有public的类属性，使用时需要用new
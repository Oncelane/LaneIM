    // 192.168.199.17 替换成你的canal server的地址
    // example 替换成-e canal.destinations=example 你自己定义的名字
    //  该字段名字在 canal\conf\example\meta.dat 文件中，NewSimpleCanalConnector函数参数配置，也在文件中
    /**
      NewSimpleCanalConnector 参数说明
        client.NewSimpleCanalConnector("Canal服务端地址", "Canal服务端端口", "Canal服务端用户名", "Canal服务端密码", "Canal服务端destination", 60000, 60*60*1000)
        Canal服务端地址：canal服务搭建地址IP
        Canal服务端端口：canal\conf\canal.properties文件中
        Canal服务端用户名、密码：canal\conf\example\instance.properties 文件中
        Canal服务端destination ：canal\conf\example\meta.dat 文件中
    */

    // https://github.com/alibaba/canal/wiki/AdminGuide
    //mysql 数据解析关注的表，Perl正则表达式.
    //
    //多个正则之间以逗号(,)分隔，转义符需要双斜杠(\\)
    //
    //常见例子：
    //
    //  1.  所有表：.*   or  .*\\..*
    //  2.  canal schema下所有表： canal\\..*
    //  3.  canal下的以canal打头的表：canal\\.canal.*
    //  4.  canal schema下的一张表：canal\\.test1
    //  5.  多个规则组合使用：canal\\..*,mysql.test1,mysql.test2 (逗号分隔)

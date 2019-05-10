alice init
bob init
charly init
alice testfilter messagesView 1 <bob> test
alice testfilter messagesView 1 <charly> test
bob testfilter messagesView 1 <alice> test
bob testfilter messagesView 1 <charly> test
charly testfilter messagesView 1 <alice> test
charly testfilter messagesView 1 <bob> test
alice msg test
bob msg test
charly msg test

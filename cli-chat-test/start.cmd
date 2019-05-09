alice init
bob init
charly init
alice testfilter messagesView 2 <bob> test
alice testfilter messagesView 1 <charly> test
bob testfilter messagesView 2 <alice> test
bob testfilter messagesView 1 <charly> test
charly testfilter messagesView 2 <alice> test
charly testfilter messagesView 2 <bob> test
alice msg test
bob msg test
charly msg test
alice msg test
bob msg test

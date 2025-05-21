import docNode from './documentNode.ts'
import makeFileNode from './makeFileNode.ts'
import sqlNode from "./sqlNode.ts";

export default {
    ...docNode,
    ...makeFileNode,
    ...sqlNode
}
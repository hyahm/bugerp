<template>
  <div class="createPost-container">
    <el-form
      ref="postForm"
      :model="postForm"
      :rules="rules"
      class="form-container"
    >
      <sticky :class-name="'sub-navbar ' + postForm.status">
        <!--<CommentDropdown v-model="postForm.comment_disabled" />-->
        <!--<PlatformDropdown v-model="postForm.platforms" />-->
        <!--<SourceUrlDropdown v-model="postForm.source_uri" />-->
        <el-button
          :disabled="ispub"
          style="margin-left: 10px;"
          type="success"
          @click="submitForm"
          >发布
        </el-button>
      </sticky>

      <div class="" style="padding: 20px">
        <el-row>
          <el-form-item label="任务类型: ">
            <el-radio-group v-model="typ" @change="handleChangeType">
              <el-radio v-for="(t, key) in ts" :key="key" :label="key">{{
                t
              }}</el-radio>
            </el-radio-group>
          </el-form-item>
          <!--<el-col :span="24" >-->
          <el-form-item prop="title" label="文章标题：">
            <el-input
              v-model="postForm.title"
              :maxlength="100"
              placeholder="请输入标题"
              clearable
              style="width: 80%;"
            />
          </el-form-item>
        </el-row>

        <el-form-item
          style="display: inline-block;width: 400px"
          label="项目名称："
        >
          <el-select
            v-model="postForm.projectname"
            placeholder="请选择"
            @change="changeProject(postForm.projectname)"
          >
            <el-option
              v-for="(item, index) in projectnames"
              :key="index"
              :label="item"
              :value="item"
            />
          </el-select>
        </el-form-item>

        <el-form-item
          v-if="typ == 1"
          style="display: inline-block;width: 400px"
          label="运行环境："
        >
          <el-select v-model="postForm.envname" placeholder="请选择">
            <el-option
              v-for="(item, index) in envnames"
              :key="index"
              :label="item"
              :value="item"
            />
          </el-select>
        </el-form-item>

        <el-form-item
          v-if="typ == 1"
          style="display: inline-block;width: 400px"
          label="应用版本："
        >
          <el-select v-model="postForm.version" placeholder="请选择">
            <el-option
              v-for="(item, index) in versions"
              :key="index"
              :label="item"
              :value="item"
            />
          </el-select>
        </el-form-item>

        <el-form-item
          v-if="typ == 1"
          style="display: inline-block;width: 400px"
          label="优先级别："
        >
          <el-select v-model="postForm.level" placeholder="请选择">
            <el-option
              v-for="(item, index) in levels"
              :key="index"
              :label="item"
              :value="item"
            />
          </el-select>
        </el-form-item>

        <!--<el-form-item style="margin-bottom: 40px;" prop="title">-->
        <!--<PlatformDropdown v-model="postForm.platforms" />-->
        <!--</el-form-item>-->
        <el-form-item
          v-if="typ == 1"
          style="display: inline-block;width: 400px"
          label="重要性："
        >
          <el-select v-model="postForm.important" placeholder="请选择">
            <el-option
              v-for="(important, index) in importants"
              :key="index"
              :label="important"
              :value="important"
            />
          </el-select>
        </el-form-item>

        <el-form-item
          style="display: inline-block;width: 400px"
          label="分配任务："
        >
          <el-select
            v-model="postForm.selectuser"
            multiple
            placeholder="分配任务"
          >
            <el-option
              v-for="(item, index) in users"
              :key="index"
              :label="item"
              :value="item"
            />
          </el-select>
        </el-form-item>

        <!--</template>-->

        <el-form-item>
          <div id="main">
            <mavon-editor
              ref="md"
              v-model="postForm.content"
              @imgAdd="imgAdd"
            />
          </div>
          <!-- <div class="editor-container">
            <Tinymce ref="editor" v-model="postForm.content" />
          </div> -->
        </el-form-item>
      </div>
    </el-form>
  </div>
</template>

<script>
// import Tinymce from '@/components/Tinymce'

import Sticky from "@/components/Sticky"; // 粘性header组件
import { fetchBug, createBug } from "@/api/bugs";
import { uploadImg } from "@/api/uploadimg";
import {
  getEnv,
  getMyProject,
  getLevels,
  getImportants,
  getProjectUser,
  getTyp
} from "@/api/get";

const defaultForm = {
  // status: 'draft',
  title: "", // 文章题目
  content: "", // 文章内容
  id: -1,
  selectuser: [],
  projectname: "",
  level: "",
  envname: "",
  important: "",
  version: "",
  typ: 1
};

export default {
  name: "ArticleDetail",
  components: {
    // Tinymce,
    Sticky
  },
  props: {
    isEdit: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      postForm: Object.assign({}, defaultForm),
      ispub: false,
      versions: [],
      importants: [],
      levels: [],
      oses: [],
      users: [],
      projectnames: [],
      envnames: [],
      ts: [],
      typ: "1"
    };
  },
  activated() {
    this.getimportants();
    this.getproject();
    this.getenv();
    this.getlevels();
  },
  created() {
    this.getimportants();
    this.getproject();
    // this.getversion()
    this.getlevels();
    this.getenv();
    this.gettyp();
    if (this.isEdit) {
      const id = this.$route.params && this.$route.params.id;
      this.postForm.id = parseInt(id);
      this.fetchData(id);
    } else {
      this.postForm = Object.assign({}, defaultForm);
    }
  },
  methods: {
    gettyp() {
      getTyp().then(resp => {
        this.ts = resp.data.ts;
      });
    },
    changeProject(name) {
      // 选择不用项目会显示不同的用户
      this.users = [];
      getProjectUser(name).then(resp => {
        this.users = resp.data.name;
        this.versions = resp.data.versions;
        this.postForm.version = "";
        this.postForm.selectuser = [];
        //
      });
    },

    editProject(name) {
      // 选择不用项目会显示不同的用户
      this.users = [];
      getProjectUser(name).then(resp => {
        this.users = resp.data.name;
        this.versions = resp.data.versions;
      });
    },
    getimportants() {
      getImportants().then(resp => {
        this.importants = resp.data.importants;
      });
    },
    getlevels() {
      getLevels().then(resp => {
        this.levels = resp.data.levels;
      });
    },
    getenv() {
      getEnv().then(resp => {
        this.envnames = resp.data.envlist;
      });
    },
    getproject() {
      getMyProject().then(resp => {
        this.projectnames = resp.data.name;
      });
    },

    fetchData(id) {
      fetchBug(id).then(resp => {
        const dd = resp.data;
        this.editProject(dd.projectname);
        this.postForm.title = dd.title;
        this.postForm.content = dd.content;
        this.postForm.important = dd.important;
        this.postForm.version = dd.version;
        this.postForm.selectuser = dd.selectuser;
        this.postForm.envname = dd.envname;
        this.postForm.level = dd.level;
        this.postForm.projectname = dd.projectname;
        this.typ = dd.typ + "";
        // 第一次获取数据， 请求到当前项目关联的版本和用户
      });
    },
    submitForm() {
      this.ispub = true;
      // this.postForm.display_time = parseInt(this.display_time / 1000)
      if (this.postForm.title.length > 40) {
        this.$message({
          message: "标题长度必须小于40位",
          type: "error"
        });
        this.ispub = false;
        return;
      }
      if (this.postForm.selectuser.length < 1) {
        this.$message({
          message: "请选择指定给谁",
          type: "error"
        });
        this.ispub = false;
        return;
      }
      if (this.postForm.projectname.length < 1) {
        this.$message({
          message: "请选择项目名称",
          type: "error"
        });
        this.ispub = false;
        return;
      }
      if (this.typ === 1 && this.postForm.level.length < 1) {
        this.$message({
          message: "请选择项目级别",
          type: "error"
        });
        this.ispub = false;
        return;
      }
      if (this.typ === 1 && this.postForm.important.length < 1) {
        this.$message({
          message: "请选择项目严重程度",
          type: "error"
        });
        this.ispub = false;
        return;
      }
      if (this.postForm.content.length < 1) {
        this.$message({
          message: "请填写内容",
          type: "error"
        });
        this.ispub = false;
        return;
      }
      if (this.typ === 1 && this.postForm.envname.length < 1) {
        this.$message({
          message: "请选择运行环境",
          type: "error"
        });
        this.ispub = false;
        return;
      }
      if (this.typ === 1 && this.postForm.version.length < 1) {
        this.$message({
          message: "请选择版本",
          type: "error"
        });
        this.ispub = false;
        return;
      }
      this.$refs.postForm.validate(valid => {
        if (valid) {
          this.postForm.typ = parseInt(this.typ);
          createBug(this.postForm).then(resp => {
            if (this.postForm.id === -1) {
              this.$notify({
                title: "成功",
                message: "发布成功",
                type: "success"
              });
              this.$router.push({ path: "/bug/edit/" + resp.data.id });
            } else {
              this.$notify({
                title: "成功",
                message: "修改成功",
                type: "success"
              });
            }
          });
        }
        this.ispub = false;
      });
    },
    imgAdd(pos, $file) {
      // 第一步.将图片上传到服务器.
      var formdata = new FormData();
      formdata.append("image", $file);
      uploadImg(formdata).then(resp => {
        this.$refs.md.$img2Url(pos, resp.data.url);
      });
    },
    draftForm() {
      if (
        this.postForm.content.length === 0 ||
        this.postForm.title.length === 0
      ) {
        this.$message({
          message: "请填写必要的标题和内容",
          type: "warning"
        });
        return;
      }
      this.$message({
        message: "保存成功",
        type: "success",
        showClose: true,
        duration: 1000
      });
    }
  }
};
</script>

<style rel="stylesheet/scss" lang="scss" scoped>
@import "src/styles/mixin.scss";

.createPost-container {
  position: relative;
  .createPost-main-container {
    padding: 40px 45px 20px 50px;
    .postInfo-container {
      position: relative;
      @include clearfix;
      margin-bottom: 10px;
      .postInfo-container-item {
        float: left;
      }
    }
    .editor-container {
      min-height: 500px;
      margin: 0 0 30px;
      .editor-upload-btn-container {
        text-align: right;
        margin-right: 10px;
        .editor-upload-btn {
          display: inline-block;
        }
      }
    }
  }
  .word-counter {
    width: 40px;
    position: absolute;
    right: -10px;
    top: 0px;
  }
}
</style>

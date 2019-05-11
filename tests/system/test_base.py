from storagebeat import BaseTest

import os


class Test(BaseTest):

    def test_base(self):
        """
        Basic test with exiting Storagebeat normally
        """
        self.render_config_template(
            path=os.path.abspath(self.working_dir) + "/log/*"
        )

        storagebeat_proc = self.start_beat()
        self.wait_until(lambda: self.log_contains("storagebeat is running"))
        exit_code = storagebeat_proc.kill_and_wait()
        assert exit_code == 0
